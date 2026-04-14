package dao

import (
	"context"
	"errors"
	"time"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type BehaviorDao interface {
	RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error
	UpdateProgress(ctx context.Context, studentID uint, progress domain.WatchProgress, status string) error
	GetClassVideoProgress(ctx context.Context, studentID uint, classID uint) ([]domain.VideoProgress, error)
}

type behaviorDao struct {
	db *gorm.DB
}

func NewBehaviorDao(db *gorm.DB) BehaviorDao {
	return &behaviorDao{db: db}
}

func (bd *behaviorDao) RecordAction(ctx context.Context, studentID uint, action domain.WatchAction) error {
	event := action.Event
	if event != "pause" && event != "replay" {
		return errors.New("unsupported event")
	}

	return bd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//先找学生行为的model(没有就创建)
		var row model.StudentBehavior
		if err := tx.Where("student_id=? AND video_id=? AND segment_id=?", studentID, action.VideoID, action.SegmentID).
			First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				row = model.StudentBehavior{
					StudentID: studentID,
					VideoID:   action.VideoID,
					SegmentID: action.SegmentID,
				}
				if err := tx.Create(&row).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		//更新行为数
		updates := map[string]interface{}{}
		if event == "pause" {
			updates["pause_count"] = gorm.Expr("pause_count + ?", 1)
			if action.Duration > 0 {
				updates["pause_duration"] = gorm.Expr("pause_duration + ?", float64(action.Duration))
			}
		}
		if event == "replay" {
			updates["replay_count"] = gorm.Expr("replay_count + ?", 1)
		}
		// 记录观看时间：pause 不计入 watch_duration
		if action.Duration > 0 && event != "pause" {
			updates["watch_duration"] = gorm.Expr("watch_duration + ?", float64(action.Duration))
		}
		if len(updates) == 0 {
			return nil
		}
		err := tx.Model(&model.StudentBehavior{}).
			Where("student_id=? AND video_id=? AND segment_id=?", studentID, action.VideoID, action.SegmentID).
			Updates(updates).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (bd *behaviorDao) UpdateProgress(ctx context.Context, studentID uint, progress domain.WatchProgress, status string) error {
	now := time.Now()
	return bd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p model.StudentVideoProgress
		err := tx.Where("user_id=? AND video_id=?", studentID, progress.VideoID).First(&p).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			row := model.StudentVideoProgress{
				UserID:        studentID,
				VideoID:       progress.VideoID,
				MaxSec:        progress.CurrentSec,
				LastSec:       progress.CurrentSec,
				Status:        status,
				WatchDuration: 0,
			}
			return tx.Create(&row).Error
		}
		if err != nil {
			return err
		}

		inc := int(now.Sub(p.UpdatedAt).Seconds())
		if inc < 0 {
			inc = 0
		}
		if inc > 600 {
			inc = 600
		}
		return tx.Model(&model.StudentVideoProgress{}).
			Where("user_id=? AND video_id=?", studentID, progress.VideoID).
			Updates(map[string]interface{}{
				"max_sec":        gorm.Expr("GREATEST(max_sec, ?)", progress.CurrentSec),
				"last_sec":       progress.CurrentSec,
				"watch_duration": gorm.Expr("watch_duration + ?", inc),
				"status":         status,
			}).Error
	})
}

func (bd *behaviorDao) GetClassVideoProgress(ctx context.Context, studentID uint, classID uint) ([]domain.VideoProgress, error) {
	// join class videos + videos + student's progress
	type row struct {
		VideoID       uint
		Title         string
		Duration      int
		Status        *string
		MaxSec        *int
		WatchDuration *int
	}
	var rows []row
	err := bd.db.WithContext(ctx).
		Table("video_to_class vtc").
		Select(`vtc.video_id as video_id,
		        v.title as title,
		        v.duration as duration,
		        p.status as status,
		        p.max_sec as max_sec,
		        p.watch_duration as watch_duration`).
		Joins("join videos v on v.id = vtc.video_id").
		Joins("left join student_video_progresses p on p.video_id = vtc.video_id and p.user_id = ?", studentID).
		Where("vtc.class_id = ?", classID).
		Order("vtc.id asc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]domain.VideoProgress, 0, len(rows))
	for _, r := range rows {
		st := "todo"
		if r.Status != nil && *r.Status != "" {
			st = *r.Status
		}
		maxSec := 0
		if r.MaxSec != nil {
			maxSec = *r.MaxSec
		}
		watchDur := 0
		if r.WatchDuration != nil {
			watchDur = *r.WatchDuration
		}
		percent := 0
		if r.Duration > 0 {
			percent = int(float64(maxSec) * 100 / float64(r.Duration))
		}
		out = append(out, domain.VideoProgress{
			VideoID:         r.VideoID,
			Title:           r.Title,
			Status:          st,
			ProgressPercent: percent,
			WatchTime:       watchDur,
		})
	}
	return out, nil
}
