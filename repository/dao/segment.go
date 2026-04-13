package dao

import (
	"context"
	"fmt"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SegmentDao interface {
	BatchUpsertSegments(context.Context, []domain.Segment) error
	UpsertKnowledgeSegmentMappings(context.Context, []domain.Segment) error
}

type segmentDao struct {
	db *gorm.DB
}

func NewSegmentDao(db *gorm.DB) SegmentDao {
	return &segmentDao{
		db: db,
	}
}

func (sd *segmentDao) BatchUpsertSegments(ctx context.Context, segments []domain.Segment) error {
	if len(segments) == 0 {
		return nil
	}

	records := make([]model.Segment, 0, len(segments))
	for _, seg := range segments {
		records = append(records, model.Segment{
			SegmentID: seg.SegmentID,
			VideoID:   seg.VideoID,
			Start:     seg.Start,
			End:       seg.End,
			Text:      seg.Text,
		})
	}

	return sd.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "segment_id"},
			},
			OnConstraint: "idx_segment_id",
			DoUpdates: clause.AssignmentColumns([]string{
				"video_id",
				"start",
				"end",
				"text",
			}),
		}).
		CreateInBatches(&records, 200).
		Error
}

func (sd *segmentDao) UpsertKnowledgeSegmentMappings(ctx context.Context, segments []domain.Segment) error {
	if len(segments) == 0 {
		return nil
	}

	segmentIDs := make([]string, 0, len(segments))
	//这里用map去重，当作set使用
	knowledgeSet := make(map[string]struct{}, 64)
	for _, seg := range segments {
		if seg.SegmentID != "" {
			segmentIDs = append(segmentIDs, seg.SegmentID)
		}
		knowledgeSet[seg.KnowledgeID] = struct{}{}
	}

	knowledgeIDs := make([]string, 0, len(knowledgeSet))
	for kid := range knowledgeSet {
		knowledgeIDs = append(knowledgeIDs, kid)
	}

	return sd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Resolve segments (string ID -> uint PK)
		var segModels []model.Segment
		if err := tx.Where("segment_id IN ?", segmentIDs).Find(&segModels).Error; err != nil {
			return err
		}
		segIDToPK := make(map[string]uint, len(segModels))
		segPKs := make([]uint, 0, len(segModels))
		for _, s := range segModels {
			segIDToPK[s.SegmentID] = s.ID
			segPKs = append(segPKs, s.ID)
		}
		if len(segIDToPK) == 0 {
			return fmt.Errorf("no segments found for mapping")
		}

		// Resolve knowledge points (string ID -> uint PK)
		var kpModels []model.KnowledgePoint
		if len(knowledgeIDs) > 0 {
			if err := tx.Where("knowledge_id IN ?", knowledgeIDs).Find(&kpModels).Error; err != nil {
				return err
			}
		}
		kpIDToPK := make(map[string]uint, len(kpModels))
		for _, kp := range kpModels {
			kpIDToPK[kp.KnowledgeID] = kp.ID
		}

		// Clear old mappings for these segments (idempotent replace)
		if err := tx.Where("segment_id IN ?", segPKs).Delete(&model.KnowledgeSegmentMapping{}).Error; err != nil {
			return err
		}

		// Insert new mappings
		mappings := make([]model.KnowledgeSegmentMapping, 0, len(segments)*2)
		for _, seg := range segments {
			segPK, ok := segIDToPK[seg.SegmentID]
			if !ok {
				return fmt.Errorf("segment not found for segment_id=%s", seg.SegmentID)
			}
			kpPK, ok := kpIDToPK[seg.KnowledgeID]
			if !ok {
				return fmt.Errorf("segment not found for knowledge_id=%s", seg.KnowledgeID)
			}
			mappings = append(mappings, model.KnowledgeSegmentMapping{
				KnowledgeID: kpPK,
				SegmentID:   segPK,
			})
		}

		if len(mappings) == 0 {
			return nil
		}
		return tx.CreateInBatches(&mappings, 500).Error
	})
}
