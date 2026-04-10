package model

import (
	"gorm.io/gorm"
)

// 记录学生视频观看进度
type StudentVideoProgress struct {
	gorm.Model
	UserID        uint   `gorm:"not null;uniqueIndex:idx_user_video;comment:学生ID"`
	VideoID       uint   `gorm:"not null;uniqueIndex:idx_user_video;comment:视频ID"`
	MaxSec        int    `gorm:"default:0;comment:最大观看秒数"`
	Status        string `gorm:"type:enum('finished','in_progress','todo','expired');default 'todo';not null;comment:观看状态"`
	WatchDuration int    `gorm:"default:0;comment:总观看时长(秒)"`
}

func (StudentVideoProgress) TableName() string {
	return "student_video_progresses"
}

type StudentBehavior struct {
	gorm.Model
	StudentID     uint    `gorm:"not null;uniqueIndex:idx_student_video_segment;comment:学生ID"`
	VideoID       uint    `gorm:"not null;uniqueIndex:idx_student_video_segment;comment:视频ID"`
	SegmentID     uint    `gorm:"not null;uniqueIndex:idx_student_video_segment;comment:分段ID"`
	PauseCount    int     `gorm:"default:0;comment:暂停次数"`
	ReplayCount   int     `gorm:"default:0;comment:回放当前分段次数"`
	WatchDuration float64 `gorm:"default:0;comment:本段观看时长"`
	PauseDuration float64 `gorm:"default:0;comment:本段暂停时长"`
}

func (StudentBehavior) TableName() string {
	return "student_behaviors"
}
