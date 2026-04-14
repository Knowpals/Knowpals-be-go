package dao

import (
	"context"
	"encoding/json"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type QuestionDao interface {
	// ReplaceQuestionsForVideo 删除该视频下旧题目及题目-知识点关联后批量写入新题目（流水线幂等）
	ReplaceQuestionsForVideo(ctx context.Context, videoID uint, items []domain.QuizQuestion) error
}

type questionDao struct {
	db *gorm.DB
}

func NewQuestionDao(db *gorm.DB) QuestionDao {
	return &questionDao{db: db}
}

func (d *questionDao) ReplaceQuestionsForVideo(ctx context.Context, videoID uint, items []domain.QuizQuestion) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldIDs []uint
		if err := tx.Model(&model.Question{}).Where("video_id = ?", videoID).Pluck("id", &oldIDs).Error; err != nil {
			return err
		}
		if len(oldIDs) > 0 {
			if err := tx.Where("question_id IN ?", oldIDs).Delete(&model.QuestionToKnowledge{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", oldIDs).Delete(&model.Question{}).Error; err != nil {
				return err
			}
		}

		for _, it := range items {

			var optsPtr *string
			if len(it.Options) > 0 {
				b, _ := json.Marshal(it.Options)
				s := string(b)
				optsPtr = &s
			}
			analysis := it.Analysis
			if analysis == "" {
				analysis = "无"
			}
			q := model.Question{
				VideoID:   videoID,
				SegmentID: it.SegmentID,
				Type:      it.Type,
				Content:   it.Content,
				Options:   optsPtr,
				Answer:    it.Answer,
				Analysis:  analysis,
			}
			if err := tx.Create(&q).Error; err != nil {
				return err
			}
			for _, kp := range it.KnowledgePKs {
				if kp == 0 {
					continue
				}
				if err := tx.Create(&model.QuestionToKnowledge{
					QuestionID:  q.ID,
					KnowledgeID: kp,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
