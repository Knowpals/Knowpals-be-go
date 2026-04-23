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
	ReplaceQuestionsForVideo(ctx context.Context, videoID uint, items []domain.Question) error
	// ReplacePersonalizedQuestionsForStudentVideo 个性化题：按 (student_id, video_id) 覆盖写入
	ReplacePersonalizedQuestionsForStudentVideo(ctx context.Context, studentID uint, videoID uint, items []domain.Question) error
	// FindKnowledgeIDsByQuestionIDs 用于写入记忆：question_id -> knowledge_id(string)（取第一条）
	FindKnowledgeIDsByQuestionIDs(ctx context.Context, questionIDs []uint) (map[uint]string, error)
	// Review: list/update/delete/add pipeline questions for teacher review
	ListPipelineQuestionsByVideo(ctx context.Context, videoID uint) ([]model.Question, error)
	UpdatePipelineQuestion(ctx context.Context, questionID uint, patch map[string]interface{}) error
	DeletePipelineQuestion(ctx context.Context, questionID uint) error
	CreatePipelineQuestion(ctx context.Context, videoID uint, segmentID *uint, q model.Question, knowledgePKs []uint) (uint, error)
	ReplaceQuestionKnowledge(ctx context.Context, questionID uint, knowledgePKs []uint) error
	// student answering
	ListQuestionsByIDsAndVideo(ctx context.Context, videoID uint, questionIDs []uint) ([]model.Question, error)
	InsertStudentAnswers(ctx context.Context, answers []model.StudentAnswer) error
	// query
	ListQuestionsByVideo(ctx context.Context, videoID uint) ([]model.Question, error)
}

type questionDao struct {
	db *gorm.DB
}

func NewQuestionDao(db *gorm.DB) QuestionDao {
	return &questionDao{db: db}
}

func (d *questionDao) ReplaceQuestionsForVideo(ctx context.Context, videoID uint, items []domain.Question) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldIDs []uint
		if err := tx.Model(&model.Question{}).Where("video_id = ? AND student_id IS NULL AND source='pipeline'", videoID).Pluck("id", &oldIDs).Error; err != nil {
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
				VideoID:    videoID,
				StudentID:  nil,
				Source:     "pipeline",
				SegmentID:  it.SegmentID,
				Type:       it.Type,
				Content:    it.Content,
				Options:    optsPtr,
				Answer:     it.Answer,
				Analysis:   analysis,
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

func (d *questionDao) ReplacePersonalizedQuestionsForStudentVideo(ctx context.Context, studentID uint, videoID uint, items []domain.Question) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldIDs []uint
		if err := tx.Model(&model.Question{}).
			Where("video_id = ? AND student_id = ? AND source='agent'", videoID, studentID).
			Pluck("id", &oldIDs).Error; err != nil {
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

		sid := studentID
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
				StudentID: &sid,
				Source:    "agent",
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

func (d *questionDao) FindKnowledgeIDsByQuestionIDs(ctx context.Context, questionIDs []uint) (map[uint]string, error) {
	out := make(map[uint]string, len(questionIDs))
	if len(questionIDs) == 0 {
		return out, nil
	}
	type row struct {
		QuestionID  uint   `gorm:"column:question_id"`
		KnowledgeID string `gorm:"column:knowledge_id"`
	}
	var rows []row
	if err := d.db.WithContext(ctx).
		Table("question_to_knowledge qk").
		Select("qk.question_id as question_id, kp.knowledge_id as knowledge_id").
		Joins("join knowledge_points kp on kp.id = qk.knowledge_id").
		Where("qk.question_id IN ?", questionIDs).
		Order("qk.id asc").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.QuestionID == 0 || r.KnowledgeID == "" {
			continue
		}
		if _, exists := out[r.QuestionID]; exists {
			continue
		}
		out[r.QuestionID] = r.KnowledgeID
	}
	return out, nil
}

func (d *questionDao) ListPipelineQuestionsByVideo(ctx context.Context, videoID uint) ([]model.Question, error) {
	var rows []model.Question
	err := d.db.WithContext(ctx).
		Where("video_id=? AND student_id IS NULL AND source='pipeline'", videoID).
		Order("id asc").
		Find(&rows).Error
	return rows, err
}

func (d *questionDao) UpdatePipelineQuestion(ctx context.Context, questionID uint, patch map[string]interface{}) error {
	return d.db.WithContext(ctx).
		Model(&model.Question{}).
		Where("id=? AND student_id IS NULL AND source='pipeline'", questionID).
		Updates(patch).Error
}

func (d *questionDao) DeletePipelineQuestion(ctx context.Context, questionID uint) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var q model.Question
		if err := tx.WithContext(ctx).
			Where("id=? AND student_id IS NULL AND source='pipeline'", questionID).
			First(&q).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).
			Where("question_id = ?", questionID).
			Delete(&model.QuestionToKnowledge{}).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Delete(&model.Question{}, questionID).Error
	})
}

func (d *questionDao) CreatePipelineQuestion(ctx context.Context, videoID uint, segmentID *uint, q model.Question, knowledgePKs []uint) (uint, error) {
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		q.ID = 0
		q.VideoID = videoID
		q.StudentID = nil
		q.Source = "pipeline"
		q.SegmentID = segmentID
		if err := tx.Create(&q).Error; err != nil {
			return err
		}
		for _, kp := range knowledgePKs {
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
		return nil
	})
	if err != nil {
		return 0, err
	}
	return q.ID, nil
}

func (d *questionDao) ReplaceQuestionKnowledge(ctx context.Context, questionID uint, knowledgePKs []uint) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// ensure question is editable (pipeline public question)
		var q model.Question
		if err := tx.Where("id=? AND student_id IS NULL AND source='pipeline'", questionID).First(&q).Error; err != nil {
			return err
		}
		if err := tx.Where("question_id = ?", questionID).Delete(&model.QuestionToKnowledge{}).Error; err != nil {
			return err
		}
		for _, kp := range knowledgePKs {
			if kp == 0 {
				continue
			}
			if err := tx.Create(&model.QuestionToKnowledge{QuestionID: questionID, KnowledgeID: kp}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *questionDao) ListQuestionsByIDsAndVideo(ctx context.Context, videoID uint, questionIDs []uint) ([]model.Question, error) {
	if len(questionIDs) == 0 {
		return []model.Question{}, nil
	}
	var rows []model.Question
	err := d.db.WithContext(ctx).
		Where("id IN ? AND video_id = ?", questionIDs, videoID).
		Find(&rows).Error
	return rows, err
}

func (d *questionDao) InsertStudentAnswers(ctx context.Context, answers []model.StudentAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).CreateInBatches(&answers, 200).Error
}

func (d *questionDao) ListQuestionsByVideo(ctx context.Context, videoID uint) ([]model.Question, error) {
	var qs []model.Question
	err := d.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("id asc").
		Find(&qs).Error
	return qs, err
}
