package model

import "gorm.io/gorm"

type KnowledgePoint struct {
	gorm.Model
	VideoID uint   `gorm:"column:video_id;type:bigint;not null;index"`
	Title   string `gorm:"column:title;type:varchar(128);not null;index"`
	Content string `gorm:"column:content;type:text;not null;index:,class:FULLTEXT"`
}

func (KnowledgePoint) TableName() string {
	return "knowledge_points"
}

// KnowledgeSegmentMapping 知识点和分段映射（可以导航）
type KnowledgeSegmentMapping struct {
	gorm.Model
	KnowledgeID uint `gorm:"column:knowledge_id;type:bigint;not null;index"`
	SegmentID   uint `gorm:"column:segment_id;type:bigint;not null;index"`
}

func (KnowledgeSegmentMapping) TableName() string {
	return "knowledge_segment_mappings"
}
