package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Knowpals/Knowpals-be-go/api/message"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/events/producer"
	"github.com/Knowpals/Knowpals-be-go/events/topic"
	"github.com/Knowpals/Knowpals-be-go/infra/cos"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/tool"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	totalPipelineStages = 2
	stageKnowledgeNum   = 1
	stageQuizNum        = 2
)

type PipelineService interface {
	CreateJob(ctx context.Context, videoID uint) (jobID string, err error)
	ProcessResult(ctx context.Context, msg *message.ResultMessage) error
}

type pipelineService struct {
	pipelineDao  dao.PipelineDao
	videoDao     dao.VideoDao
	knowledgeDao dao.KnowledgeDao
	segmentDao   dao.SegmentDao
	producer     producer.Producer
	cos          *cos.COSClient
	log          *zap.Logger
}

func NewPipelineService(
	pdao dao.PipelineDao,
	vdao dao.VideoDao,
	kdao dao.KnowledgeDao,
	sdao dao.SegmentDao,
	prod producer.Producer,
	cos *cos.COSClient,
	log *zap.Logger,
) PipelineService {
	return &pipelineService{
		pipelineDao:  pdao,
		videoDao:     vdao,
		knowledgeDao: kdao,
		segmentDao:   sdao,
		producer:     prod,
		cos:          cos,
		log:          log,
	}
}

func (s *pipelineService) CreateJob(ctx context.Context, videoID uint) (string, error) {
	if s.producer == nil {
		return "", errors.CreateJobError(fmt.Errorf("kafka producer 未初始化"))
	}
	jobID := uuid.New().String()
	job := &domain.PipelineJob{
		JobID:        jobID,
		VideoID:      videoID,
		Status:       "running",
		CurrentStage: stageKnowledgeNum,
		TotalStage:   totalPipelineStages,
	}
	//创建任务
	if err := s.pipelineDao.CreateJob(ctx, job); err != nil {
		return "", errors.CreateJobError(err)
	}

	jobStage := &domain.PipelineJobStage{
		JobID:      jobID,
		Stage:      message.StageKnowledge,
		Status:     "running",
		RetryCount: 0,
	}
	err := s.pipelineDao.CreateStage(ctx, jobStage)
	if err != nil {
		//把创建的任务也失败
		_ = s.pipelineDao.UpdateJob(ctx, jobID, map[string]interface{}{
			"status": "failed",
		})
		return "", errors.CreateJobError(err)
	}

	video, err := s.videoDao.GetVideoByID(ctx, videoID)
	if err != nil {
		_ = s.failJob(ctx, jobID, message.StageKnowledge, err.Error())
		return "", errors.CreateJobError(err)
	}

	url, err := s.cos.SignUrl(ctx, video.FileKey)
	if err != nil {
		_ = s.failJob(ctx, jobID, message.StageKnowledge, err.Error())
		return "", errors.CreateJobError(err)
	}

	payload := map[string]interface{}{"file_urls": []string{url}}
	task := message.TaskMessage{
		JobID:   jobID,
		Stage:   message.StageKnowledge,
		Payload: payload,
		Retry:   0,
	}
	if err = s.producer.SendTask(topic.TASK_TOPIC, task); err != nil {
		_ = s.failJob(ctx, jobID, message.StageKnowledge, err.Error())
		return "", errors.CreateJobError(err)
	}

	return jobID, nil
}

func (s *pipelineService) ProcessResult(ctx context.Context, msg *message.ResultMessage) error {
	switch msg.Stage {
	case message.StageKnowledge:
		return s.runKnowledgeStage(ctx, msg)
	case message.StageQuiz:
		return s.runQuizStage(ctx, msg)
	default:
		return errors.ProcessTaskError(fmt.Errorf("unknown stage: %s", msg.Stage))
	}
}

func (s *pipelineService) runKnowledgeStage(ctx context.Context, msg *message.ResultMessage) error {

	s.log.Info(fmt.Sprintf("[jobID:%s] 进入 knowledge stage"))

	if msg.Status == "failed" || msg.Error != "" {
		//把状态更新为失败
		_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, msg.Error)
		return errors.RunKnowledgeStageError(fmt.Errorf("执行知识点拆分任务失败：%s", msg.Error))
	}

	job, err := s.pipelineDao.GetJob(ctx, msg.JobID)
	if err != nil {
		_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
		return errors.RunKnowledgeStageError(err)
	}
	videoID := job.VideoID

	if err := s.ensureStage(ctx, msg.JobID, message.StageKnowledge); err != nil {
		_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
		return errors.RunKnowledgeStageError(err)
	}

	var result message.KnowledgeSegmentPayload
	err = tool.MapToStruct(msg.Result, &result)
	if err != nil {
		_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
		return errors.RunKnowledgeStageError(err)
	}

	//更新video的duration
	err = s.videoDao.UpdateVideo(ctx, videoID, map[string]interface{}{"duration": result.DurationMs})
	if err != nil {
		_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
		return errors.RunKnowledgeStageError(err)
	}

	// concepts -> knowledge_points
	kps := make([]domain.KnowledgePoint, 0, len(result.Concepts))
	for _, c := range result.Concepts {
		kps = append(kps, domain.KnowledgePoint{
			KnowledgeID: c.ConceptID,
			VideoID:     videoID,
			Title:       c.Title,
			Content:     c.Content,
		})
	}
	if s.knowledgeDao != nil {
		if err := s.knowledgeDao.BatchUpsertKnowledgePoints(ctx, kps); err != nil {
			_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
			return errors.RunKnowledgeStageError(err)
		}
	}

	// segments -> segments (+ text) and mappings
	segs := make([]domain.Segment, 0, len(result.Segments))
	for _, seg := range result.Segments {
		segs = append(segs, domain.Segment{
			SegmentID:   seg.SegmentID,
			VideoID:     videoID,
			Start:       seg.StartMs,
			End:         seg.EndMs,
			Text:        seg.Text,
			KnowledgeID: seg.ConceptID,
		})
	}
	if s.segmentDao != nil {
		if err := s.segmentDao.BatchUpsertSegments(ctx, segs); err != nil {
			_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
			return errors.RunKnowledgeStageError(err)
		}
		if err := s.segmentDao.UpsertKnowledgeSegmentMappings(ctx, segs); err != nil {
			_ = s.failJob(ctx, msg.JobID, message.StageKnowledge, err.Error())
			return errors.RunKnowledgeStageError(err)
		}
	}

	out, _ := json.Marshal(msg.Result)
	if err := s.pipelineDao.UpdateStage(ctx, msg.JobID, message.StageKnowledge, map[string]interface{}{
		"status": "success",
		"output": string(out),
	}); err != nil {
		s.log.Error(fmt.Sprintf("[jobID:%s stage:%s] update stage失败：%v", msg.JobID, message.StageKnowledge, err))
		return errors.RunKnowledgeStageError(err)
	}
	if err := s.pipelineDao.UpdateJob(ctx, msg.JobID, map[string]interface{}{
		"current_stage": stageQuizNum,
	}); err != nil {
		s.log.Error(fmt.Sprintf("[jobID:%s stage:%s] update job失败：%v", msg.JobID, message.StageKnowledge, err))
		return errors.RunKnowledgeStageError(err)
	}

	s.log.Info(fmt.Sprintf("[jobID:%s] 完成 knowledge stage"))

	//创建新阶段
	jobStage := &domain.PipelineJobStage{
		JobID:      msg.JobID,
		Stage:      message.StageQuiz,
		Status:     "running",
		RetryCount: 0,
	}
	err = s.pipelineDao.CreateStage(ctx, jobStage)
	if err != nil {
		s.log.Error(fmt.Sprintf("更新stage错误：%v", err))
		return errors.RunKnowledgeStageError(err)
	}

	next := message.TaskMessage{
		JobID:   msg.JobID,
		Stage:   message.StageQuiz,
		Payload: map[string]interface{}{"video_id": videoID},
		Retry:   0,
	}
	if s.producer != nil {
		err := s.producer.SendTask(topic.TASK_TOPIC, next)
		if err != nil {
			s.log.Error(fmt.Sprintf("kafka:send next stage task message error:%v", err))
			return errors.RunQuizStageError(err)
		}
	}
	return nil
}

func (s *pipelineService) runQuizStage(ctx context.Context, msg *message.ResultMessage) error {

	s.log.Info(fmt.Sprintf("[jobID:%s] 进入 quiz stage"))

	if msg.Status == "failed" || msg.Error != "" {
		_ = s.failJob(ctx, msg.JobID, message.StageQuiz, msg.Error)
		return errors.RunQuizStageError(fmt.Errorf(msg.Error))
	}
	if err := s.ensureStage(ctx, msg.JobID, message.StageQuiz); err != nil {
		_ = s.failJob(ctx, msg.JobID, message.StageQuiz, msg.Error)
		return errors.RunQuizStageError(err)
	}

	out, _ := json.Marshal(msg.Result)
	if err := s.pipelineDao.UpdateStage(ctx, msg.JobID, message.StageQuiz, map[string]interface{}{
		"status": "success",
		"output": string(out),
	}); err != nil {
		s.log.Error(fmt.Sprintf("[jobID:%s stage:%s] update stage 失败：%v", msg.JobID, message.StageKnowledge, err))
		return errors.RunQuizStageError(err)
	}
	if err := s.pipelineDao.UpdateJob(ctx, msg.JobID, map[string]interface{}{
		"status":        "success",
		"current_stage": stageQuizNum,
	}); err != nil {
		s.log.Error(fmt.Sprintf("[jobID:%s stage:%s] update job 失败：%v", msg.JobID, message.StageKnowledge, err))
		return errors.RunQuizStageError(err)
	}

	s.log.Info(fmt.Sprintf("[jobID:%s] 完成 quiz stage"))

	return nil
}

func (s *pipelineService) failJob(ctx context.Context, jobID string, stage string, errMsg string) error {
	_ = s.ensureStage(ctx, jobID, stage)
	_ = s.pipelineDao.UpdateStage(ctx, jobID, stage, map[string]interface{}{"status": "failed", "output": errMsg})
	_ = s.pipelineDao.UpdateJob(ctx, jobID, map[string]interface{}{"status": "failed"})
	return fmt.Errorf("%s", errMsg)
}

func (s *pipelineService) ensureStage(ctx context.Context, jobID string, stage string) error {
	n, err := s.pipelineDao.CheckStage(ctx, jobID, stage)
	if err != nil {
		return err
	}
	if n == 0 {
		return s.pipelineDao.CreateStage(ctx, &domain.PipelineJobStage{
			JobID:  jobID,
			Stage:  stage,
			Status: "running",
		})
	}
	return s.pipelineDao.UpdateStage(ctx, jobID, stage, map[string]interface{}{
		"status": "running",
	})
}
