package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VideoDao interface {
	SaveVideo(ctx context.Context, video domain.Video) (uint, error)
	GetVideoByID(ctx context.Context, id uint) (domain.Video, error)
	UpdateVideo(ctx context.Context, id uint, updates map[string]interface{}) error
	ListSegmentsByVideoID(ctx context.Context, videoID uint) ([]domain.Segment, error)
	ListKnowledgePointsByVideoID(ctx context.Context, videoID uint) ([]domain.KnowledgePoint, error)
	ListQuestionsByVideoID(ctx context.Context, videoID uint) ([]domain.Question, error)
	ListQuestionKnowledge(ctx context.Context, questionIDs []uint) (map[uint][]domain.KnowledgePoint, error)
	AssignVideoToClasses(ctx context.Context, videoID uint, classIDs []uint) error
	ListClassVideoTasks(ctx context.Context, classID uint) ([]domain.Video, error)
}

type videoDao struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewVideoDao(db *gorm.DB, log *zap.Logger) VideoDao {
	return &videoDao{db: db, log: log}
}

func (vd *videoDao) SaveVideo(ctx context.Context, video domain.Video) (uint, error) {
	videoModel := model.Video{
		TeacherID: video.TeacherID,
		Title:     video.Title,
		FileKey:   video.FileKey,
	}
	if err := vd.db.WithContext(ctx).Create(&videoModel).Error; err != nil {
		return 0, err
	}
	return videoModel.ID, nil
}

func (vd *videoDao) GetVideoByID(ctx context.Context, id uint) (domain.Video, error) {
	var m model.Video
	if err := vd.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return domain.Video{}, err
	}
	return domain.Video{
		ID:        m.ID,
		TeacherID: m.TeacherID,
		FileKey:   m.FileKey,
		Title:     m.Title,
		Duration:  m.Duration,
	}, nil
}

func (vd *videoDao) UpdateVideo(ctx context.Context, id uint, updates map[string]interface{}) error {
	return vd.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", id).Updates(updates).Error
}

func (vd *videoDao) ListSegmentsByVideoID(ctx context.Context, videoID uint) ([]domain.Segment, error) {
	var segs []model.Segment
	if err := vd.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("start asc").
		Find(&segs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Segment, 0, len(segs))
	for _, s := range segs {
		out = append(out, domain.Segment{
			ID:        s.ID,
			SegmentID: s.SegmentID,
			VideoID:   s.VideoID,
			Start:     s.Start,
			End:       s.End,
			Text:      s.Text,
		})
	}
	return out, nil
}

func (vd *videoDao) ListKnowledgePointsByVideoID(ctx context.Context, videoID uint) ([]domain.KnowledgePoint, error) {
	var kps []model.KnowledgePoint
	if err := vd.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("id asc").
		Find(&kps).Error; err != nil {
		return nil, err
	}
	out := make([]domain.KnowledgePoint, 0, len(kps))
	for _, kp := range kps {
		out = append(out, domain.KnowledgePoint{
			ID:          kp.ID,
			KnowledgeID: kp.KnowledgeID,
			VideoID:     kp.VideoID,
			Title:       kp.Title,
			Content:     kp.Content,
		})
	}
	return out, nil
}

func (vd *videoDao) ListQuestionsByVideoID(ctx context.Context, videoID uint) ([]domain.Question, error) {
	var qsModel []model.Question
	if err := vd.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("id asc").
		Find(&qsModel).Error; err != nil {
		return nil, err
	}

	qs := make([]domain.Question, len(qsModel))
	for i, q := range qsModel {
		if q.Options != nil {
			var op []string
			err := json.Unmarshal([]byte(*q.Options), &op)
			if err != nil {
				vd.log.Error(fmt.Sprintf("option格式解析错误:%v", err))
			}
			qs[i] = domain.Question{
				ID:        q.ID,
				SegmentID: q.SegmentID,
				Type:      q.Type,
				Content:   q.Content,
				Options:   op,
				Answer:    q.Answer,
				Analysis:  q.Analysis,
			}
		} else {
			qs[i] = domain.Question{
				ID:        q.ID,
				SegmentID: q.SegmentID,
				Type:      q.Type,
				Content:   q.Content,
				Answer:    q.Answer,
				Analysis:  q.Analysis,
			}
		}

	}

	return qs, nil
}

func (vd *videoDao) ListQuestionKnowledge(ctx context.Context, questionIDs []uint) (map[uint][]domain.KnowledgePoint, error) {
	out := make(map[uint][]domain.KnowledgePoint, len(questionIDs))
	if len(questionIDs) == 0 {
		return out, nil
	}
	var links []model.QuestionToKnowledge
	if err := vd.db.WithContext(ctx).
		Where("question_id IN ?", questionIDs).
		Find(&links).Error; err != nil {
		return nil, err
	}
	knowledgeIDs := make([]uint, 0, len(links))
	qidToKIDs := make(map[uint][]uint, len(questionIDs))
	for _, l := range links {
		qidToKIDs[l.QuestionID] = append(qidToKIDs[l.QuestionID], l.KnowledgeID)
		knowledgeIDs = append(knowledgeIDs, l.KnowledgeID)
	}
	var kps []model.KnowledgePoint
	if len(knowledgeIDs) > 0 {
		if err := vd.db.WithContext(ctx).Where("id IN ?", knowledgeIDs).Find(&kps).Error; err != nil {
			return nil, err
		}
	}
	kpByID := make(map[uint]model.KnowledgePoint, len(kps))
	for _, kp := range kps {
		kpByID[kp.ID] = kp
	}
	for qid, kids := range qidToKIDs {
		for _, kid := range kids {
			kp, ok := kpByID[kid]
			if !ok {
				continue
			}
			out[qid] = append(out[qid], domain.KnowledgePoint{
				ID:          kp.ID,
				KnowledgeID: kp.KnowledgeID,
				VideoID:     kp.VideoID,
				Title:       kp.Title,
				Content:     kp.Content,
			})
		}
	}
	return out, nil
}

func (vd *videoDao) AssignVideoToClasses(ctx context.Context, videoID uint, classIDs []uint) error {
	if len(classIDs) == 0 {
		return nil
	}
	return vd.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 幂等：先删后插
		if err := tx.Where("video_id = ? AND class_id IN ?", videoID, classIDs).Delete(&model.VideoToClass{}).Error; err != nil {
			return err
		}
		rows := make([]model.VideoToClass, 0, len(classIDs))
		for _, cid := range classIDs {
			rows = append(rows, model.VideoToClass{VideoID: videoID, ClassID: cid})
		}
		return tx.CreateInBatches(&rows, 200).Error
	})
}

func (vd *videoDao) ListClassVideoTasks(ctx context.Context, classID uint) ([]domain.Video, error) {
	var videoIDs []uint
	if err := vd.db.WithContext(ctx).
		Model(&model.VideoToClass{}).
		Where("class_id = ?", classID).
		Pluck("video_id", &videoIDs).Error; err != nil {
		return nil, err
	}
	if len(videoIDs) == 0 {
		return []domain.Video{}, nil
	}
	var vs []model.Video
	if err := vd.db.WithContext(ctx).Where("id IN ?", videoIDs).Find(&vs).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Video, 0, len(vs))
	for _, v := range vs {
		out = append(out, domain.Video{
			ID:        v.ID,
			TeacherID: v.TeacherID,
			FileKey:   v.FileKey,
			Title:     v.Title,
			Duration:  v.Duration,
			CreatedAt: v.CreatedAt,
		})
	}
	return out, nil
}
