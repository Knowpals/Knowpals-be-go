package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KnowledgeDao interface {
	BatchUpsertKnowledgePoints(context.Context, []domain.KnowledgePoint) error
	FindKnowledgePointsByKnowledgeIDs(ctx context.Context, knowledgeIDs []string) (map[string]domain.KnowledgePoint, error)
}

type knowledgeDao struct {
	db *gorm.DB
}

func NewKnowledgeDao(db *gorm.DB) KnowledgeDao {
	return &knowledgeDao{
		db: db,
	}
}

func (kd *knowledgeDao) BatchUpsertKnowledgePoints(ctx context.Context, kps []domain.KnowledgePoint) error {
	if len(kps) == 0 {
		return nil
	}

	records := make([]model.KnowledgePoint, 0, len(kps))
	for _, kp := range kps {
		records = append(records, model.KnowledgePoint{
			KnowledgeID: kp.KnowledgeID,
			VideoID:     kp.VideoID,
			Title:       kp.Title,
			Content:     kp.Content,
		})
	}

	return kd.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "knowledge_id"},
			},
			OnConstraint: "idx_knowledge_id",
			DoUpdates: clause.AssignmentColumns([]string{
				"video_id",
				"title",
				"content",
			}),
		}).
		CreateInBatches(&records, 200).
		Error
}

func (kd *knowledgeDao) FindKnowledgePointsByKnowledgeIDs(ctx context.Context, knowledgeIDs []string) (map[string]domain.KnowledgePoint, error) {
	out := make(map[string]domain.KnowledgePoint, len(knowledgeIDs))
	if len(knowledgeIDs) == 0 {
		return out, nil
	}

	var records []model.KnowledgePoint
	if err := kd.db.WithContext(ctx).
		Where("knowledge_id IN ?", knowledgeIDs).
		Find(&records).Error; err != nil {
		return nil, err
	}

	for _, r := range records {
		out[r.KnowledgeID] = domain.KnowledgePoint{
			ID:          r.ID,
			KnowledgeID: r.KnowledgeID,
			VideoID:     r.VideoID,
			Title:       r.Title,
			Content:     r.Content,
		}
	}
	return out, nil
}
