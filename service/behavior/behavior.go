package behavior

import (
	"context"
	"time"

	"github.com/Knowpals/Knowpals-be-go/domain"
	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/service/agentclient"
)

type BehaviorService interface {
	RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error
	UpdateProgress(ctx context.Context, studentID uint, progress domain.WatchProgress) error
	GetClassVideoProgress(ctx context.Context, studentID uint, classID uint, status string) ([]domain.VideoProgress, error)
	GetMyUnfinishedTasks(ctx context.Context, studentID uint) ([]domain.UnfinishedTask, error)
}

type behaviorService struct {
	behaviorDao dao.BehaviorDao
	videoDao    dao.VideoDao
	mem         agentclient.MemoryWriter
	segmentDao  dao.SegmentDao
}

func NewBehaviorService(behaviorDao dao.BehaviorDao, videoDao dao.VideoDao, segmentDao dao.SegmentDao, mem agentclient.MemoryWriter) BehaviorService {
	return &behaviorService{behaviorDao: behaviorDao, videoDao: videoDao, segmentDao: segmentDao, mem: mem}
}

func (bs *behaviorService) RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error {
	err := bs.behaviorDao.RecordAction(ctx, studentID, action)
	if err != nil {
		return errors2.RecordActionError(err)
	}
	if bs.mem != nil {
		// best-effort：写入短期记忆失败不影响主流程
		kid := ""
		if bs.segmentDao != nil && action.SegmentID != 0 {
			kid, _ = bs.segmentDao.GetKnowledgeIDBySegmentPK(ctx, action.SegmentID)
		}
		if kid != "" {
			if action.Event == "pause" {
				_ = bs.mem.WritePause(ctx, studentID, kid, action.VideoID, action.SegmentID)
			}
			if action.Event == "replay" {
				_ = bs.mem.WriteReplay(ctx, studentID, kid, action.VideoID, action.SegmentID)
			}
		}
	}
	return nil
}

func (bs *behaviorService) UpdateProgress(ctx context.Context, studentID uint, progress domain.WatchProgress) error {
	video, err := bs.videoDao.GetVideoByID(ctx, progress.VideoID)
	if err != nil {
		return errors2.UpdateProgressError(err)
	}
	status := "todo"
	if video.Duration > 0 && progress.CurrentSec >= int(float64(video.Duration)*0.9) && time.Now().Before(video.Deadline) {
		status = "finished"
	}

	if progress.CurrentSec < int(float64(video.Duration)*0.9) && time.Now().After(video.Deadline) {
		status = "expired"
	}
	if err := bs.behaviorDao.UpdateProgress(ctx, studentID, progress, status); err != nil {
		return errors2.UpdateProgressError(err)
	}
	return nil
}

func (bs *behaviorService) GetClassVideoProgress(ctx context.Context, studentID uint, classID uint, status string) ([]domain.VideoProgress, error) {
	out, err := bs.behaviorDao.GetClassVideoProgress(ctx, studentID, classID, status)
	if err != nil {
		return nil, errors2.GetClassVideoProgressError(err)
	}
	return out, nil
}

func (bs *behaviorService) GetMyUnfinishedTasks(ctx context.Context, studentID uint) ([]domain.UnfinishedTask, error) {
	out, err := bs.behaviorDao.ListMyUnfinishedTasks(ctx, studentID)
	if err != nil {
		return nil, errors2.GetMyUnfinishedTasksError(err)
	}
	return out, nil
}
