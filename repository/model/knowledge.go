package model

import "gorm.io/gorm"

type KnowledgePoint struct {
	gorm.Model
	KnowledgeID string `gorm:"column:knowledge_id;type:varchar(64);uniqueIndex"`
	VideoID     uint   `gorm:"column:video_id;not null;index"`
	Title       string `gorm:"column:title;type:varchar(128);not null;index"`
	Content     string `gorm:"column:content;type:text;not null;index:,class:FULLTEXT"`
}

func (KnowledgePoint) TableName() string {
	return "knowledge_points"
}

// KnowledgeSegmentMapping 知识点和分段映射
// 注意：字段名不要用 SegmentID / KnowledgePointID 这类 GORM 默认关联命名，否则 Create 时可能被当成 BelongsTo 外键处理，写入 0。
type KnowledgeSegmentMapping struct {
	gorm.Model
	KnowledgePk uint `gorm:"column:knowledge_id;not null;index"`
	SegmentPk   uint `gorm:"column:segment_id;not null;index"`
}

func (KnowledgeSegmentMapping) TableName() string {
	return "knowledge_segment_mappings"
}
