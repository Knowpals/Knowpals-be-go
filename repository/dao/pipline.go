package dao

import (
	"context"
	"errors"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type PipelineDao interface {
	CreateJob(ctx context.Context, job *domain.PipelineJob) error
	GetJob(ctx context.Context, jobID string) (*domain.PipelineJob, error)
	UpdateJob(ctx context.Context, jobID string, updates map[string]interface{}) error
	CreateStage(ctx context.Context, stage *domain.PipelineJobStage) error
	UpdateStage(ctx context.Context, jobID string, stage string, updates map[string]interface{}) error
	CheckStage(ctx context.Context, jobID string, stage string) (int64, error)
	ListStages(ctx context.Context, jobID string) ([]domain.PipelineJobStage, error)
}

type pipelineDao struct {
	db *gorm.DB
}

func NewPipelineDao(db *gorm.DB) PipelineDao {
	return &pipelineDao{db: db}
}

func (d *pipelineDao) CreateJob(ctx context.Context, job *domain.PipelineJob) error {
	m := &model.PipelineJob{
		JobID:        job.JobID,
		VideoID:      job.VideoID,
		Status:       job.Status,
		CurrentStage: job.CurrentStage,
		TotalStage:   job.TotalStage,
	}
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *pipelineDao) GetJob(ctx context.Context, jobID string) (*domain.PipelineJob, error) {
	var j model.PipelineJob
	err := d.db.WithContext(ctx).Where("job_id = ?", jobID).First(&j).Error
	if err != nil {
		return nil, err
	}
	return &domain.PipelineJob{
		JobID:        jobID,
		VideoID:      j.VideoID,
		Status:       j.Status,
		CurrentStage: j.CurrentStage,
		TotalStage:   j.TotalStage,
	}, nil
}

func (d *pipelineDao) UpdateJob(ctx context.Context, jobID string, updates map[string]interface{}) error {
	res := d.db.WithContext(ctx).Model(&model.PipelineJob{}).Where("job_id = ?", jobID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("pipeline job not found")
	}
	return nil
}

func (d *pipelineDao) CreateStage(ctx context.Context, stage *domain.PipelineJobStage) error {
	m := &model.PipelineJobStage{
		JobID:      stage.JobID,
		Stage:      stage.Stage,
		Status:     stage.Status,
		RetryCount: stage.RetryCount,
	}
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *pipelineDao) UpdateStage(ctx context.Context, jobID string, stage string, updates map[string]interface{}) error {
	res := d.db.WithContext(ctx).Model(&model.PipelineJobStage{}).
		Where("job_id = ? AND stage = ?", jobID, stage).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("pipeline stage not found")
	}
	return nil
}

func (d *pipelineDao) CheckStage(ctx context.Context, jobID string, stage string) (int64, error) {
	var n int64
	err := d.db.WithContext(ctx).Model(&model.PipelineJobStage{}).
		Where("job_id = ? AND stage = ?", jobID, stage).Count(&n).Error
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (d *pipelineDao) ListStages(ctx context.Context, jobID string) ([]domain.PipelineJobStage, error) {
	var stages []model.PipelineJobStage
	if err := d.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("id asc").
		Find(&stages).Error; err != nil {
		return nil, err
	}
	out := make([]domain.PipelineJobStage, 0, len(stages))
	for _, s := range stages {
		out = append(out, domain.PipelineJobStage{
			ID:         s.ID,
			JobID:      s.JobID,
			Stage:      s.Stage,
			Status:     s.Status,
			RetryCount: s.RetryCount,
			Output:     s.Output,
		})
	}
	return out, nil
}
