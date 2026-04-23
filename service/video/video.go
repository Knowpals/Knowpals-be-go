package video

import (
	"context"
	errors1 "errors"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
)

type VideoService interface {
	SaveVideo(ctx context.Context, video domain.Video) (uint, error)
	GetVideo(ctx context.Context, videoID uint) (domain.Video, error)
	GetVideoDetail(ctx context.Context, videoID uint) (domain.Video, []domain.Segment, []domain.KnowledgePoint, []domain.Question, map[uint][]domain.KnowledgePoint, error)
	StartReview(ctx context.Context, videoID uint) error
	Publish(ctx context.Context, videoID uint) error
	AssignVideoToClasses(ctx context.Context, videoID uint, classIDs []uint) error
	ListClassVideoTasks(ctx context.Context, classID uint) ([]domain.Video, error)
	ListMyUploadedVideos(ctx context.Context, teacherID uint) ([]domain.Video, error)
}

type videoService struct {
	dao dao.VideoDao
}

func NewVideoService(dao dao.VideoDao) VideoService {
	return &videoService{dao: dao}
}

func (vs *videoService) SaveVideo(ctx context.Context, video domain.Video) (uint, error) {
	return vs.dao.SaveVideo(ctx, video)
}

func (vs *videoService) GetVideo(ctx context.Context, videoID uint) (domain.Video, error) {
	return vs.dao.GetVideoByID(ctx, videoID)
}

func (vs *videoService) GetVideoDetail(ctx context.Context, videoID uint) (domain.Video, []domain.Segment, []domain.KnowledgePoint, []domain.Question, map[uint][]domain.KnowledgePoint, error) {
	v, err := vs.dao.GetVideoByID(ctx, videoID)
	if err != nil {
		return domain.Video{}, nil, nil, nil, nil, errors.UploadVideoError(err)
	}
	segs, err := vs.dao.ListSegmentsByVideoID(ctx, videoID)
	if err != nil {
		return domain.Video{}, nil, nil, nil, nil, errors.UploadVideoError(err)
	}
	kps, err := vs.dao.ListKnowledgePointsByVideoID(ctx, videoID)
	if err != nil {
		return domain.Video{}, nil, nil, nil, nil, errors.UploadVideoError(err)
	}
	qs, err := vs.dao.ListQuestionsByVideoID(ctx, videoID)
	if err != nil {
		return domain.Video{}, nil, nil, nil, nil, err
	}
	qids := make([]uint, 0, len(qs))
	for _, q := range qs {
		qids = append(qids, q.ID)
	}
	qkps, err := vs.dao.ListQuestionKnowledge(ctx, qids)
	if err != nil {
		return domain.Video{}, nil, nil, nil, nil, errors.UploadVideoError(err)
	}
	return v, segs, kps, qs, qkps, nil
}

func (vs *videoService) AssignVideoToClasses(ctx context.Context, videoID uint, classIDs []uint) error {
	st, err := vs.dao.GetVideoReviewStatus(ctx, videoID)
	if err != nil {
		return errors.UploadVideoError(err)
	}
	if st != "published" {
		return errors.VideoNotPublishedError(errors1.New("未发布不可下发"))
	}
	return vs.dao.AssignVideoToClasses(ctx, videoID, classIDs)
}

func (vs *videoService) ListClassVideoTasks(ctx context.Context, classID uint) ([]domain.Video, error) {
	return vs.dao.ListClassVideoTasks(ctx, classID)
}

func (vs *videoService) ListMyUploadedVideos(ctx context.Context, teacherID uint) ([]domain.Video, error) {
	return vs.dao.ListVideosByTeacher(ctx, teacherID)
}
