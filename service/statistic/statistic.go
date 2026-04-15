package statistic

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	errors2 "github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
)

type StatService interface {
	GetStudentStat(ctx context.Context, studentID uint, videoID uint) (domain.StudentVideoStat, error)
	GetClassStat(ctx context.Context, classID uint, videoID uint) (domain.ClassVideoStat, error)
	GetStudentOverview(ctx context.Context, studentID uint) (domain.StudentOverviewStat, error)
}

type statService struct {
	dao dao.StatisticDao
}

func NewStatService(dao dao.StatisticDao) StatService {
	return &statService{dao: dao}
}

func (s *statService) GetStudentStat(ctx context.Context, studentID uint, videoID uint) (domain.StudentVideoStat, error) {
	stat, err := s.dao.GetStudentVideoStat(ctx, studentID, videoID)
	if err != nil {
		return domain.StudentVideoStat{}, errors2.GetStudentStatError(err)
	}
	return stat, nil
}

func (s *statService) GetClassStat(ctx context.Context, classID uint, videoID uint) (domain.ClassVideoStat, error) {
	stat, err := s.dao.GetClassVideoStat(ctx, classID, videoID)
	if err != nil {
		return domain.ClassVideoStat{}, errors2.GetClassStatError(err)
	}

	return stat, nil
}

func (s *statService) GetStudentOverview(ctx context.Context, studentID uint) (domain.StudentOverviewStat, error) {
	stat, err := s.dao.GetStudentOverviewStat(ctx, studentID)
	if err != nil {
		return domain.StudentOverviewStat{}, errors2.GetStudentOverviewError(err)
	}
	return stat, nil
}
