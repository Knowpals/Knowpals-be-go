package model

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	VideoID   uint    `gorm:"column:video_id;type:bigint;not null;index"`
	SegmentID *uint   `gorm:"column:segment_id;type:bigint;index"` //为空为整视频题，非空为分段题
	Type      string  `gorm:"column:type;type:enum('choice','fill','judge');not null"`
	Content   string  `gorm:"column:content;type:text;not null"`
	Options   *string `gorm:"column:options;type:text"`
	Answer    string  `gorm:"column:answer;type:text;not null"`
	Analysis  string  `gorm:"column:analysis;type:text;not null"`
}

func (Question) TableName() string {
	return "questions"
}

// QuestionToKnowledge 题目和知识点映射表
type QuestionToKnowledge struct {
	gorm.Model
	QuestionID  uint `gorm:"column:question_id;type:bigint;not null;index"`
	KnowledgeID uint `gorm:"column:knowledge_id;type:bigint;not null;index"`
}

func (QuestionToKnowledge) TableName() string {
	return "question_to_knowledge"
}

type StudentAnswer struct {
	gorm.Model
	StudentID  uint    `gorm:"not null;index:idx_student_question;index:idx_student_video;comment:学生ID"`
	QuestionID uint    `gorm:"not null;index:idx_student_question;comment:题目ID"`
	IsCorrect  bool    `gorm:"not null;comment:是否正确"`
	TimeCost   float64 `gorm:"comment:答题耗时(秒)"`
	VideoID    uint    `gorm:"not null;index:idx_student_video;comment:冗余-视频ID，方便统计"`
	SegmentID  *uint   `gorm:"index;comment:冗余-分段ID，可为null（整体题）"`
}

func (StudentAnswer) TableName() string {
	return "student_answers"
}
