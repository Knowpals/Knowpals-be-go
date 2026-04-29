package model

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Video{},
		&Class{},

		&Segment{},
		&KnowledgePoint{},
		&Question{},

		&ClassStudent{},
		&VideoToClass{},

		&PipelineJob{},
		&PipelineJobStage{},
		&KnowledgeSegmentMapping{},
		&QuestionToKnowledge{},
		&StudentAnswer{},
		&StudentBehavior{},
		&StudentVideoProgress{},
		&ChatMessage{},
		&StudentReport{},
	)

}
