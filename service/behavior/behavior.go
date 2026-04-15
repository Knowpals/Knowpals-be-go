package behavior

import (
	"context"
	"time"

	"github.com/Knowpals/Knowpals-be-go/domain"
	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
)

type BehaviorService interface {
	RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error
	UpdateProgress(ctx context.Context, studentID uint, progress domain.WatchProgress) error
	GetClassVideoProgress(ctx context.Context, studentID uint, classID uint, status string) ([]domain.VideoProgress, error)
}

type behaviorService struct {
	behaviorDao dao.BehaviorDao
	videoDao    dao.VideoDao
}

func NewBehaviorService(behaviorDao dao.BehaviorDao, videoDao dao.VideoDao) BehaviorService {
	return &behaviorService{behaviorDao: behaviorDao, videoDao: videoDao}
}

func (bs *behaviorService) RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error {
	err := bs.behaviorDao.RecordAction(ctx, studentID, action)
	if err != nil {
		return errors2.RecordActionError(err)
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
