package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Knowpals/Knowpals-be-go/api/grpc/agentpb"
	httpAgent "github.com/Knowpals/Knowpals-be-go/api/http/agent"
	"github.com/Knowpals/Knowpals-be-go/domain"
	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type AgentService interface {
	Chat(ctx context.Context, studentID uint, req httpAgent.ChatReq) (httpAgent.ChatResp, error)
	GenerateQuiz(ctx context.Context, studentID uint, req httpAgent.GenerateQuizReq) (httpAgent.GenerateQuizResp, error)
	GenerateReport(ctx context.Context, studentID uint, req httpAgent.GenerateReportReq) (httpAgent.GenerateReportResp, error)
	GetChatHistory(ctx context.Context, studentID uint, req httpAgent.GetChatHistoryReq) (httpAgent.GetChatHistoryResp, error)
}

type agentService struct {
	client       agentpb.AgentServiceClient
	knowledgeDao dao.KnowledgeDao
	questionDao  dao.QuestionDao
	chatDao      dao.ChatDao
	reportDao    dao.ReportDao
}

func NewAgentService(client agentpb.AgentServiceClient, knowledgeDao dao.KnowledgeDao, questionDao dao.QuestionDao, chatDao dao.ChatDao, reportDao dao.ReportDao) AgentService {
	return &agentService{client: client, knowledgeDao: knowledgeDao, questionDao: questionDao, chatDao: chatDao, reportDao: reportDao}
}

func (s *agentService) Chat(ctx context.Context, studentID uint, req httpAgent.ChatReq) (httpAgent.ChatResp, error) {
	resp, err := s.client.Chat(ctx, &agentpb.ChatRequest{
		StudentId:   strconv.Itoa(int(studentID)),
		Text:        req.Text,
		VideoId:     req.VideoID,
		KnowledgeId: req.KnowledgeID,
	})
	if err != nil {
		return httpAgent.ChatResp{}, errors2.AgentChatError(err)
	}
	out := httpAgent.ChatResp{Reply: resp.Reply, Context: resp.Context, VideoID: resp.VideoId}
	if s.chatDao != nil {
		// best-effort：存储失败不影响对话返回
		_ = s.saveChatPair(ctx, studentID, req, out.Reply)
	}
	return out, nil
}

func (s *agentService) GenerateQuiz(ctx context.Context, studentID uint, req httpAgent.GenerateQuizReq) (httpAgent.GenerateQuizResp, error) {
	resp, err := s.client.GenerateQuiz(ctx, &agentpb.GenerateQuizRequest{
		StudentId:    strconv.Itoa(int(studentID)),
		VideoId:      req.VideoID,
		NumQuestions: req.NumQuestions,
	})
	if err != nil {
		return httpAgent.GenerateQuizResp{}, errors2.AgentGenerateQuizError(err)
	}

	vid64, err := strconv.ParseUint(req.VideoID, 10, 64)
	if err != nil {
		return httpAgent.GenerateQuizResp{}, errors2.AgentGenerateQuizError(err)
	}
	videoID := uint(vid64)

	// resolve knowledge_id(string) -> knowledge_points.id(uint)
	knowledgeIDs := make([]string, 0, len(resp.Quizzes))
	for _, q := range resp.Quizzes {
		if q.KnowledgeId != "" {
			knowledgeIDs = append(knowledgeIDs, q.KnowledgeId)
		}
	}
	kmap, err := s.knowledgeDao.FindKnowledgePointsByKnowledgeIDs(ctx, knowledgeIDs)
	if err != nil {
		return httpAgent.GenerateQuizResp{}, errors2.AgentGenerateQuizError(err)
	}

	// persist personalized questions (reuse questions/question_to_knowledge)
	items := make([]domain.Question, 0, len(resp.Quizzes))
	out := make([]httpAgent.QuizItem, 0, len(resp.Quizzes))
	for _, q := range resp.Quizzes {
		kp := kmap[q.KnowledgeId]
		if kp.ID == 0 {
			// knowledge_id 必须可解析，否则不落库也不返回（防止后续写入记忆等链路缺失）
			continue
		}
		items = append(items, domain.Question{
			SegmentID:    nil,
			Type:         q.Type,
			Content:      q.Question,
			Options:      q.Options,
			Answer:       q.Answer,
			Analysis:     q.Analysis,
			KnowledgePKs: []uint{kp.ID},
		})
		out = append(out, httpAgent.QuizItem{
			KnowledgeID: q.KnowledgeId,
			Type:        q.Type,
			Question:    q.Question,
			Options:     q.Options,
			Answer:      q.Answer,
			Analysis:    q.Analysis,
			Difficulty:  q.Difficulty,
		})
	}
	if err := s.questionDao.ReplacePersonalizedQuestionsForStudentVideo(ctx, studentID, videoID, items); err != nil {
		return httpAgent.GenerateQuizResp{}, errors2.AgentGenerateQuizError(err)
	}
	return httpAgent.GenerateQuizResp{Quizzes: out}, nil
}

func (s *agentService) GenerateReport(ctx context.Context, studentID uint, req httpAgent.GenerateReportReq) (httpAgent.GenerateReportResp, error) {
	vid64, err := strconv.ParseUint(req.VideoID, 10, 64)
	if err != nil {
		return httpAgent.GenerateReportResp{}, errors2.AgentGenerateReportError(err)
	}
	videoID := uint(vid64)

	if !req.ForceRegen && s.reportDao != nil {
		existing, err := s.reportDao.GetByStudentVideo(ctx, studentID, videoID)
		if err == nil && existing.Report != "" {
			var out httpAgent.GenerateReportResp
			if e := json.Unmarshal([]byte(existing.Report), &out); e == nil {
				return out, nil
			}
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return httpAgent.GenerateReportResp{}, errors2.AgentGenerateReportError(err)
		}
	}

	resp, err := s.client.GenerateReport(ctx, &agentpb.GenerateReportRequest{
		StudentId: strconv.Itoa(int(studentID)),
		VideoId:   req.VideoID,
	})
	if err != nil {
		return httpAgent.GenerateReportResp{}, errors2.AgentGenerateReportError(err)
	}
	items := make([]httpAgent.ReportItem, 0, len(resp.Items))
	for _, it := range resp.Items {
		segs := make([]httpAgent.SegmentRef, 0, len(it.RecommendedSegments))
		for _, sref := range it.RecommendedSegments {
			if sref == nil {
				continue
			}
			segs = append(segs, httpAgent.SegmentRef{
				VideoID:   sref.VideoId,
				SegmentID: sref.SegmentId,
				StartMs:   sref.StartMs,
				EndMs:     sref.EndMs,
			})
		}
		items = append(items, httpAgent.ReportItem{
			KnowledgeID:         it.KnowledgeId,
			Mastery:             it.Mastery,
			Summary:             it.Summary,
			Weakness:            it.Weakness,
			BehaviorPattern:     it.BehaviorPattern,
			Trend:               it.Trend,
			RecommendedSegments: segs,
		})
	}
	out := httpAgent.GenerateReportResp{
		VideoID:        resp.VideoId,
		Items:          items,
		OverallSummary: resp.OverallSummary,
	}
	if s.reportDao != nil {
		if b, e := json.Marshal(out); e == nil {
			_ = s.reportDao.Upsert(ctx, model.StudentReport{
				StudentID: studentID,
				VideoID:   videoID,
				Report:    string(b),
			})
		}
	}
	return out, nil
}

func (s *agentService) GetChatHistory(ctx context.Context, studentID uint, req httpAgent.GetChatHistoryReq) (httpAgent.GetChatHistoryResp, error) {
	if s.chatDao == nil {
		return httpAgent.GetChatHistoryResp{Messages: []httpAgent.ChatMessage{}}, nil
	}
	var vidPtr *uint
	if req.VideoID != "" {
		vid64, err := strconv.ParseUint(req.VideoID, 10, 64)
		if err != nil {
			return httpAgent.GetChatHistoryResp{}, errors2.AgentChatError(err)
		}
		vid := uint(vid64)
		vidPtr = &vid
	}
	rows, err := s.chatDao.ListHistory(ctx, studentID, vidPtr, req.Limit)
	if err != nil {
		return httpAgent.GetChatHistoryResp{}, errors2.AgentChatError(err)
	}
	out := make([]httpAgent.ChatMessage, 0, len(rows))
	for _, r := range rows {
		v := ""
		if r.VideoID != nil {
			v = strconv.FormatUint(uint64(*r.VideoID), 10)
		}
		out = append(out, httpAgent.ChatMessage{
			Role:        r.Role,
			Text:        r.Text,
			VideoID:     v,
			KnowledgeID: r.KnowledgeID,
			CreatedAt:   r.CreatedAt,
		})
	}
	return httpAgent.GetChatHistoryResp{Messages: out}, nil
}
