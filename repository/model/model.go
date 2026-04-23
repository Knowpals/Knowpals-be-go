package model

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Class{},
		&ClassStudent{},
		&ChatMessage{},
		&StudentReport{},
		&KnowledgePoint{},
		&KnowledgeSegmentMapping{},
		&Question{},
		&QuestionToKnowledge{},
		&StudentAnswer{},
		&StudentBehavior{},
		&StudentVideoProgress{},
		&Video{},
		&VideoToClass{},
		&Segment{},
		&PipelineJob{},
		&PipelineJobStage{},
	)

}
