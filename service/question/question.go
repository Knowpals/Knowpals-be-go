package question

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	httpQuestion "github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/Knowpals/Knowpals-be-go/domain"
	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"github.com/Knowpals/Knowpals-be-go/service/agentclient"
	"gorm.io/gorm"
)

type Service interface {
	AnswerQuestion(ctx context.Context, studentID uint, req httpQuestion.AnswerQuestionReq) (httpQuestion.AnswerQuestionResp, error)
	GenerateVideoExercise(ctx context.Context, studentID uint, role domain.RoleType, videoID uint) (httpQuestion.GenerateVideoExerciseResp, error)

	ReviewListQuestions(ctx context.Context, teacherID uint, role domain.RoleType, videoID uint) (httpQuestion.GenerateVideoExerciseResp, error)
	ReviewAddQuestion(ctx context.Context, teacherID uint, role domain.RoleType, req httpQuestion.ReviewAddReq) error
	ReviewUpdateQuestion(ctx context.Context, teacherID uint, role domain.RoleType, req httpQuestion.ReviewUpdateReq) error
	ReviewDeleteQuestion(ctx context.Context, teacherID uint, role domain.RoleType, questionID uint) error
}

type service struct {
	qdao dao.QuestionDao
	sdao dao.SegmentDao
	vdao dao.VideoDao
	mem  agentclient.MemoryWriter
}

func NewQuestionService(qdao dao.QuestionDao, sdao dao.SegmentDao, vdao dao.VideoDao, mem agentclient.MemoryWriter) Service {
	return &service{qdao: qdao, sdao: sdao, vdao: vdao, mem: mem}
}

func (s *service) AnswerQuestion(ctx context.Context, studentID uint, req httpQuestion.AnswerQuestionReq) (httpQuestion.AnswerQuestionResp, error) {
	if len(req.StudentAnswers) == 0 {
		return httpQuestion.AnswerQuestionResp{Results: []httpQuestion.Result{}}, nil
	}
	qids := make([]uint, 0, len(req.StudentAnswers))
	for _, a := range req.StudentAnswers {
		qids = append(qids, a.QuestionID)
	}

	questions, err := s.qdao.ListQuestionsByIDsAndVideo(ctx, req.VideoID, qids)
	if err != nil {
		return httpQuestion.AnswerQuestionResp{}, err
	}
	qmap := make(map[uint]model.Question, len(questions))
	for _, q := range questions {
		qmap[q.ID] = q
	}

	normalize := func(st string) string { return strings.TrimSpace(strings.ToLower(st)) }
	results := make([]httpQuestion.Result, 0, len(req.StudentAnswers))
	toInsert := make([]model.StudentAnswer, 0, len(req.StudentAnswers))

	for _, a := range req.StudentAnswers {
		q, ok := qmap[a.QuestionID]
		if !ok {
			continue
		}
		isCorrect := normalize(a.Answer) == normalize(q.Answer)
		results = append(results, httpQuestion.Result{
			QuestionID: q.ID,
			IsCorrect:  isCorrect,
			Answer:     q.Answer,
			Analysis:   q.Analysis,
		})
		toInsert = append(toInsert, model.StudentAnswer{
			StudentID:  studentID,
			QuestionID: q.ID,
			IsCorrect:  isCorrect,
			TimeCost:   float64(a.TimeCost),
			VideoID:    req.VideoID,
			SegmentID:  q.SegmentID,
		})
	}

	if err := s.qdao.InsertStudentAnswers(ctx, toInsert); err != nil {
		return httpQuestion.AnswerQuestionResp{}, err
	}

	// best-effort memory write (require knowledge_id)
	if s.mem != nil {
		kmap, _ := s.qdao.FindKnowledgeIDsByQuestionIDs(ctx, qids)
		for _, a := range req.StudentAnswers {
			q, ok := qmap[a.QuestionID]
			if !ok {
				continue
			}
			kid := kmap[q.ID]
			if kid == "" {
				continue
			}
			isCorrect := normalize(a.Answer) == normalize(q.Answer)
			contentIfWrong := ""
			if !isCorrect {
				contentIfWrong = q.Content
			}
			_ = s.mem.WriteQuestion(ctx, studentID, kid, req.VideoID, q.SegmentID, q.ID, isCorrect, a.Answer, q.Answer, contentIfWrong)
		}
	}

	return httpQuestion.AnswerQuestionResp{Results: results}, nil
}

func (s *service) GenerateVideoExercise(ctx context.Context, studentID uint, role domain.RoleType, videoID uint) (httpQuestion.GenerateVideoExerciseResp, error) {
	if role == domain.Role_Student {
		st, err := s.vdao.GetVideoReviewStatus(ctx, videoID)
		if err != nil {
			return httpQuestion.GenerateVideoExerciseResp{}, errors2.GetVideoDetailError(err)
		}
		if st != "published" {
			return httpQuestion.GenerateVideoExerciseResp{}, errors2.VideoNotPublishedError(errors.New("未发布"))
		}
	}

	qs, err := s.qdao.ListQuestionsByVideo(ctx, videoID)
	if err != nil {
		return httpQuestion.GenerateVideoExerciseResp{}, err
	}
	out := make([]httpQuestion.Question, 0, len(qs))
	for _, q := range qs {
		opts := []string{}
		if q.Options != nil && *q.Options != "" {
			_ = json.Unmarshal([]byte(*q.Options), &opts)
		}
		out = append(out, httpQuestion.Question{
			ID:        q.ID,
			Type:      q.Type,
			Content:   q.Content,
			Options:   opts,
			Answer:    q.Answer,
			Analysis:  q.Analysis,
			SegmentID: q.SegmentID,
		})
	}
	return httpQuestion.GenerateVideoExerciseResp{Questions: out}, nil
}

func (s *service) ReviewListQuestions(ctx context.Context, teacherID uint, role domain.RoleType, videoID uint) (httpQuestion.GenerateVideoExerciseResp, error) {
	if role != domain.Role_Teacher {
		return httpQuestion.GenerateVideoExerciseResp{}, errors.New("无权限")
	}
	qs, err := s.qdao.ListPipelineQuestionsByVideo(ctx, videoID)
	if err != nil {
		return httpQuestion.GenerateVideoExerciseResp{}, err
	}
	out := make([]httpQuestion.Question, 0, len(qs))
	for _, q := range qs {
		opts := []string{}
		if q.Options != nil && *q.Options != "" {
			_ = json.Unmarshal([]byte(*q.Options), &opts)
		}
		out = append(out, httpQuestion.Question{
			ID:        q.ID,
			Type:      q.Type,
			Content:   q.Content,
			Options:   opts,
			Answer:    q.Answer,
			Analysis:  q.Analysis,
			SegmentID: q.SegmentID,
		})
	}
	return httpQuestion.GenerateVideoExerciseResp{Questions: out}, nil
}

func (s *service) ReviewAddQuestion(ctx context.Context, teacherID uint, role domain.RoleType, req httpQuestion.ReviewAddReq) error {
	if role != domain.Role_Teacher {
		return errors.New("无权限")
	}
	segID, err := s.sdao.FindSegmentPKByVideoAndTimeMs(ctx, req.VideoID, req.TimeMs)
	if err != nil {
		return err
	}
	if segID == nil || *segID == 0 {
		return errors.New("无法定位到所属分段")
	}

	knowledgePKs, err := s.sdao.GetKnowledgePKsBySegmentPK(ctx, *segID)
	if err != nil {
		return err
	}
	if len(knowledgePKs) == 0 {
		return errors.New("该分段未关联知识点")
	}

	var optsPtr *string
	if len(req.Options) > 0 {
		b, _ := json.Marshal(req.Options)
		st := string(b)
		optsPtr = &st
	}
	analysis := req.Analysis
	if analysis == "" {
		analysis = "无"
	}
	_, err = s.qdao.CreatePipelineQuestion(ctx, req.VideoID, segID, model.Question{
		Type:     req.Type,
		Content:  req.Content,
		Options:  optsPtr,
		Answer:   req.Answer,
		Analysis: analysis,
	}, knowledgePKs)
	return err
}

func (s *service) ReviewUpdateQuestion(ctx context.Context, teacherID uint, role domain.RoleType, req httpQuestion.ReviewUpdateReq) error {
	if role != domain.Role_Teacher {
		return errors.New("无权限")
	}
	patch := map[string]interface{}{}
	if req.Type != nil {
		patch["type"] = *req.Type
	}
	if req.Content != nil {
		patch["content"] = *req.Content
	}
	if req.Answer != nil {
		patch["answer"] = *req.Answer
	}
	if req.Analysis != nil {
		patch["analysis"] = *req.Analysis
	}
	if req.Options != nil {
		b, _ := json.Marshal(req.Options)
		st := string(b)
		patch["options"] = &st
	}
	if len(patch) > 0 {
		if err := s.qdao.UpdatePipelineQuestion(ctx, req.QuestionID, patch); err != nil {
			return err
		}
	}
	if len(req.KnowledgePKs) > 0 {
		if err := s.qdao.ReplaceQuestionKnowledge(ctx, req.QuestionID, req.KnowledgePKs); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ReviewDeleteQuestion(ctx context.Context, teacherID uint, role domain.RoleType, questionID uint) error {
	if role != domain.Role_Teacher {
		return errors.New("无权限")
	}
	return s.qdao.DeletePipelineQuestion(ctx, questionID)
}

var _ *gorm.DB
