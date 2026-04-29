package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title        string     `gorm:"column:title;type:varchar(100);not null"`
	FileKey      string     `gorm:"column:file_key;type:varchar(100);not null"`
	TeacherID    uint       `gorm:"column:teacher_id;not null"`
	Duration     int        `gorm:"column:duration;not null"`
	Deadline     time.Time  `gorm:"column:deadline;type:datetime"`
	ReviewStatus string     `gorm:"column:review_status;type:enum('processing','reviewing','published');not null;default 'processing';index"`
	ReviewedAt   *time.Time `gorm:"column:reviewed_at;type:datetime"`
	PublishedAt  *time.Time `gorm:"column:published_at;type:datetime"`

	Teacher User `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (Video) TableName() string {
	return "videos"
}

// VideoToClass 班级的视频权限表
type VideoToClass struct {
	gorm.Model
	VideoID uint `gorm:"column:video_id;not null"`
	ClassID uint `gorm:"column:class_id;not null"`

	Video Video `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Class Class `gorm:"foreignKey:ClassID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (VideoToClass) TableName() string {
	return "video_to_class"
}

// Segment 视频分段
type Segment struct {
	gorm.Model
	// 业务侧的 segment string id（来自 pipeline），不要用 column=segment_id，
	// 否则会和其它表的 segment_id 外键列（引用 segments.id）冲突，导致迁移时生成反向外键。
	SegmentSID string `gorm:"column:segment_sid;type:varchar(64);uniqueIndex"`
	VideoID   uint   `gorm:"column:video_id;not null;index"`
	Start     int    `gorm:"column:start;type:int;not null"`
	End       int    `gorm:"column:end;type:int;not null"`
	Text      string `gorm:"column:text;type:text"`

	Video Video `gorm:"foreignKey:VideoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Segment) TableName() string {
	return "segments"
}
