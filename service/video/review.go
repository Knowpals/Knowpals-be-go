package video

import (
	"context"
	"errors"
	"time"

	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"gorm.io/gorm"
)

func (vs *videoService) StartReview(ctx context.Context, videoID uint) error {
	// processing -> reviewing
	st, err := vs.dao.GetVideoReviewStatus(ctx, videoID)
	if err != nil {
		return errors2.VideoReviewError(err)
	}
	if st == "" {
		st = "processing"
	}
	if st == "published" {
		return errors2.VideoReviewError(errors.New("视频已发布"))
	}
	// 允许重复调用：reviewing 直接返回成功
	if st == "reviewing" {
		return nil
	}
	now := time.Now()
	if err := vs.dao.UpdateVideo(ctx, videoID, map[string]interface{}{
		"review_status": "reviewing",
		"reviewed_at":   &now,
	}); err != nil {
		return errors2.VideoReviewError(err)
	}
	return nil
}

func (vs *videoService) Publish(ctx context.Context, videoID uint) error {
	st, err := vs.dao.GetVideoReviewStatus(ctx, videoID)
	if err != nil {
		return errors2.VideoPublishError(err)
	}
	if st != "reviewing" && st != "published" {
		return errors2.VideoPublishError(errors.New("请先进入审核"))
	}
	if st == "published" {
		return nil
	}
	now := time.Now()
	if err := vs.dao.UpdateVideo(ctx, videoID, map[string]interface{}{
		"review_status": "published",
		"published_at":  &now,
	}); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors2.VideoPublishError(err)
	}
	return nil
}

