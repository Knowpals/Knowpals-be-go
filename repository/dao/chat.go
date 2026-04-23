package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type ChatDao interface {
	Save(ctx context.Context, msg model.ChatMessage) error
	ListHistory(ctx context.Context, studentID uint, videoID *uint, limit int) ([]model.ChatMessage, error)
}

type chatDao struct {
	db *gorm.DB
}

func NewChatDao(db *gorm.DB) ChatDao {
	return &chatDao{db: db}
}

func (d *chatDao) Save(ctx context.Context, msg model.ChatMessage) error {
	return d.db.WithContext(ctx).Create(&msg).Error
}

func (d *chatDao) ListHistory(ctx context.Context, studentID uint, videoID *uint, limit int) ([]model.ChatMessage, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := d.db.WithContext(ctx).Model(&model.ChatMessage{}).Where("student_id = ?", studentID)
	if videoID != nil {
		q = q.Where("video_id = ?", *videoID)
	}
	var rows []model.ChatMessage
	if err := q.Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	// 前端一般按时间正序渲染，这里反转为 asc
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	return rows, nil
}

