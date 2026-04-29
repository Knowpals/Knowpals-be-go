package model

import (
	"time"
)

type PipelineJob struct {
	JobID        string `gorm:"column:job_id;type:varchar(64);primaryKey"`
	VideoID      uint   `gorm:"column:video_id;not null;index"`
	Status       string `gorm:"type:enum('running','success','failed');not null"`
	CurrentStage int    `gorm:"column:current_stage;type:int;not null"`
	TotalStage   int    `gorm:"column:total_stage;type:int;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Video Video `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PipelineJob) TableName() string {
	return "pipeline_jobs"
}

type PipelineJobStage struct {
	ID         uint   `gorm:"type:bigint;autoIncrement;primaryKey"`
	JobID      string `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uniq_job_stage"`
	Stage      string `gorm:"column:stage;type:varchar(32);not null;uniqueIndex:uniq_job_stage"`
	Status     string `gorm:"type:enum('running','success','failed');not null"`
	RetryCount int    `gorm:"column:retry_count;type:int;default:0"`
	Output     string `gorm:"type:text"`
	UpdatedAt  time.Time

	Job PipelineJob `gorm:"foreignKey:JobID;references:JobID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PipelineJobStage) TableName() string {
	return "pipeline_job_stages"
}
