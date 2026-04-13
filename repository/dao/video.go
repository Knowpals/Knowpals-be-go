package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type VideoDao interface {
	SaveVideo(ctx context.Context, video domain.Video) (uint, error)
	GetVideoByID(ctx context.Context, id uint) (domain.Video, error)
	UpdateVideo(ctx context.Context, id uint, updates map[string]interface{}) error
}

type videoDao struct {
	db *gorm.DB
}

func NewVideoDao(db *gorm.DB) VideoDao {
	return &videoDao{db: db}
}

func (vd *videoDao) SaveVideo(ctx context.Context, video domain.Video) (uint, error) {
	videoModel := model.Video{
		TeacherID: video.TeacherID,
		Title:     video.Title,
		FileKey:   video.FileKey,
	}
	if err := vd.db.WithContext(ctx).Create(&videoModel).Error; err != nil {
		return 0, err
	}
	return videoModel.ID, nil
}

func (vd *videoDao) GetVideoByID(ctx context.Context, id uint) (domain.Video, error) {
	var m model.Video
	if err := vd.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return domain.Video{}, err
	}
	return domain.Video{
		ID:        m.ID,
		TeacherID: m.TeacherID,
		FileKey:   m.FileKey,
		Title:     m.Title,
		Duration:  m.Duration,
	}, nil
}

func (vd *videoDao) UpdateVideo(ctx context.Context, id uint, updates map[string]interface{}) error {
	return vd.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", id).Updates(updates).Error
}
