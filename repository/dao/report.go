package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportDao interface {
	GetByStudentVideo(ctx context.Context, studentID uint, videoID uint) (model.StudentReport, error)
	Upsert(ctx context.Context, row model.StudentReport) error
}

type reportDao struct {
	db *gorm.DB
}

func NewReportDao(db *gorm.DB) ReportDao {
	return &reportDao{db: db}
}

func (d *reportDao) GetByStudentVideo(ctx context.Context, studentID uint, videoID uint) (model.StudentReport, error) {
	var r model.StudentReport
	err := d.db.WithContext(ctx).
		Where("student_id=? AND video_id=?", studentID, videoID).
		Order("id desc").
		First(&r).Error
	return r, err
}

func (d *reportDao) Upsert(ctx context.Context, row model.StudentReport) error {
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "student_id"}, {Name: "video_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"report",
				"updated_at",
			}),
		}).
		Create(&row).Error
}

