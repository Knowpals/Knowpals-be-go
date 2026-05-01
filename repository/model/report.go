package model

import "gorm.io/gorm"

// StudentReport 学习报告（个性化，按 student_id + video_id 覆盖）
type StudentReport struct {
	gorm.Model
	StudentID uint   `gorm:"column:student_id;not null;index:idx_student_video,unique"`
	VideoID   uint   `gorm:"column:video_id;not null;index:idx_student_video,unique"`
	Report    string `gorm:"column:report;type:longtext;not null"` // JSON of api/http/agent.GenerateReportResp
}

func (StudentReport) TableName() string {
	return "student_reports"
}
