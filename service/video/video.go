package video

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
)

type VideoService interface {
	SaveVideo(ctx context.Context, video domain.Video) (uint, error)
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
