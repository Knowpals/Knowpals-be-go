package model

import "gorm.io/gorm"

// ChatMessage 学生-助手对话记录（用于前端渲染历史消息）
type ChatMessage struct {
	gorm.Model
	StudentID   uint   `gorm:"column:student_id;not null;index:idx_student_video_created"`
	Role        string `gorm:"column:role;type:enum('user','assistant');not null"`
	Text        string `gorm:"column:text;type:text;not null"`
	VideoID     *uint  `gorm:"column:video_id;type:bigint;index:idx_student_video_created"`
	KnowledgeID string `gorm:"column:knowledge_id;type:varchar(64);index"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
