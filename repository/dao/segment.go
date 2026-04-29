package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type SegmentDao interface {
	BatchUpsertSegments(context.Context, []domain.Segment) error
	UpsertKnowledgeSegmentMappings(context.Context, []domain.Segment) error
	BatchGetSegmentBySegmentID(context.Context, []string) (map[string]domain.Segment, error)
	GetKnowledgeIDBySegmentPK(context.Context, uint) (string, error)
	FindSegmentPKByVideoAndTimeMs(ctx context.Context, videoID uint, timeMs int64) (*uint, error)
	GetKnowledgePKsBySegmentPK(ctx context.Context, segmentPK uint) ([]uint, error)
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

	// 不用 CreateInBatches+OnConflict：GORM 在部分 MySQL 驱动/版本下批量 upsert 可能只命中首条或错误合并为同一行。
	return sd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seg := range segments {
			if seg.SegmentID == "" {
				return fmt.Errorf("segment_id is required")
			}
			var existing model.Segment
			err := tx.Where("segment_sid = ?", seg.SegmentID).Take(&existing).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				row := model.Segment{
					SegmentSID: seg.SegmentID,
					VideoID:   seg.VideoID,
					Start:     seg.Start,
					End:       seg.End,
					Text:      seg.Text,
				}
				if err := tx.Create(&row).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"video_id": seg.VideoID,
				"start":    seg.Start,
				"end":      seg.End,
				"text":     seg.Text,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (sd *segmentDao) UpsertKnowledgeSegmentMappings(ctx context.Context, segments []domain.Segment) error {
	if len(segments) == 0 {
		return nil
	}

	segmentIDSet := make(map[string]struct{}, len(segments))
	knowledgeSet := make(map[string]struct{}, 64)
	var pairs []domain.Segment
	for _, seg := range segments {
		if seg.SegmentID == "" {
			continue
		}
		segmentIDSet[seg.SegmentID] = struct{}{}
		if seg.KnowledgeID == "" {
			continue
		}
		knowledgeSet[seg.KnowledgeID] = struct{}{}
		pairs = append(pairs, seg)
	}

	segmentIDs := make([]string, 0, len(segmentIDSet))
	for sid := range segmentIDSet {
		segmentIDs = append(segmentIDs, sid)
	}
	knowledgeIDs := make([]string, 0, len(knowledgeSet))
	for kid := range knowledgeSet {
		knowledgeIDs = append(knowledgeIDs, kid)
	}

	return sd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var segModels []model.Segment
		if err := tx.Where("segment_sid IN ?", segmentIDs).Find(&segModels).Error; err != nil {
			return err
		}
		segIDToPK := make(map[string]uint, len(segModels))
		segPKs := make([]uint, 0, len(segModels))
		for _, s := range segModels {
			segIDToPK[s.SegmentSID] = s.ID
			segPKs = append(segPKs, s.ID)
		}
		if len(segIDToPK) == 0 {
			return fmt.Errorf("no segments found for mapping")
		}

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

		if err := tx.Where("segment_id IN ?", segPKs).Delete(&model.KnowledgeSegmentMapping{}).Error; err != nil {
			return err
		}

		mappings := make([]model.KnowledgeSegmentMapping, 0, len(pairs))
		seenPair := make(map[string]struct{}, len(pairs))
		for _, seg := range pairs {
			segPK, ok := segIDToPK[seg.SegmentID]
			if !ok {
				return fmt.Errorf("segment not found for segment_id=%s", seg.SegmentID)
			}
			kpPK, ok := kpIDToPK[seg.KnowledgeID]
			if !ok {
				return fmt.Errorf("knowledge point not found for knowledge_id=%s", seg.KnowledgeID)
			}
			if segPK == 0 || kpPK == 0 {
				return fmt.Errorf("invalid pk for mapping segment_id=%s knowledge_id=%s", seg.SegmentID, seg.KnowledgeID)
			}
			key := fmt.Sprintf("%d:%d", kpPK, segPK)
			if _, dup := seenPair[key]; dup {
				continue
			}
			seenPair[key] = struct{}{}
			mappings = append(mappings, model.KnowledgeSegmentMapping{
				KnowledgePk: kpPK,
				SegmentPk:   segPK,
			})
		}

		if len(mappings) == 0 {
			return nil
		}
		return tx.CreateInBatches(&mappings, 500).Error
	})
}

func (sd *segmentDao) BatchGetSegmentBySegmentID(ctx context.Context, segmentID []string) (map[string]domain.Segment, error) {
	out := make(map[string]domain.Segment, len(segmentID))
	if len(segmentID) == 0 {
		return out, nil
	}

	var records []model.Segment
	err := sd.db.WithContext(ctx).Where("segment_sid IN ?", segmentID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		out[r.SegmentSID] = domain.Segment{
			ID:        r.ID,
			SegmentID: r.SegmentSID,
			VideoID:   r.VideoID,
			Start:     r.Start,
			End:       r.End,
			Text:      r.Text,
		}
	}

	return out, nil
}

func (sd *segmentDao) GetKnowledgeIDBySegmentPK(ctx context.Context, segmentPK uint) (string, error) {
	if segmentPK == 0 {
		return "", nil
	}
	type row struct {
		KnowledgeID string `gorm:"column:knowledge_id"`
	}
	var r row
	err := sd.db.WithContext(ctx).
		Table("knowledge_segment_mappings ksm").
		Select("kp.knowledge_id as knowledge_id").
		Joins("join knowledge_points kp on kp.id = ksm.knowledge_id").
		Where("ksm.segment_id = ?", segmentPK).
		Order("ksm.id asc").
		Limit(1).
		Scan(&r).Error
	if err != nil {
		return "", err
	}
	return r.KnowledgeID, nil
}

func (sd *segmentDao) FindSegmentPKByVideoAndTimeMs(ctx context.Context, videoID uint, timeMs int64) (*uint, error) {
	if videoID == 0 {
		return nil, nil
	}
	sec := int(timeMs / 1000)
	var seg model.Segment
	err := sd.db.WithContext(ctx).
		Where("video_id = ? AND start <= ? AND end >= ?", videoID, sec, sec).
		Order("start asc").
		First(&seg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	id := seg.ID
	return &id, nil
}

func (sd *segmentDao) GetKnowledgePKsBySegmentPK(ctx context.Context, segmentPK uint) ([]uint, error) {
	if segmentPK == 0 {
		return []uint{}, nil
	}
	var ids []uint
	if err := sd.db.WithContext(ctx).
		Model(&model.KnowledgeSegmentMapping{}).
		Where("segment_id = ?", segmentPK).
		Pluck("knowledge_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
