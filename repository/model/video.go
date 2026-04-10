package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title     string    `gorm:"column:title;type:varchar(100);not null"`
	Url       string    `gorm:"column:url;type:varchar(100);not null"`
	TeacherID uint      `gorm:"column:teacher_id;type:bigint;not null"`
	Duration  int       `gorm:"column:duration;not null"`
	Deadline  time.Time `gorm:"column:deadline;not null"`
}

func (Video) TableName() string {
	return "videos"
}

// VideoToClass 班级的视频权限表
type VideoToClass struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;not null"`
	ClassID uint `gorm:"column:class_id;not null"`
}

func (VideoToClass) TableName() string {
	return "video_to_class"
}

// Segment 视频分段
type Segment struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;type:bigint;not null;index"`
	Start   int  `gorm:"column:start;type:int;not null"`
	End     int  `gorm:"column:end;type:int;not null"`
}

func (Segment) TableName() string {
	return "segments"
}
