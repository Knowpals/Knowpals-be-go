package class

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/tool"
)

type ClassService interface {
	CreateClass(ctx context.Context, teacherID uint, className string) (domain.Class, error)
	JoinClass(ctx context.Context, userID uint, inviteCode string) error
	QuitClass(ctx context.Context, userID uint, classID uint) error
	GetStudentClasses(ctx context.Context, userID uint) ([]domain.Class, error)
	GetTeacherClasses(ctx context.Context, teacherID uint) ([]domain.Class, error)
	GetClassByID(ctx context.Context, classID uint) (domain.Class, error)
	GetClassStudents(ctx context.Context, classID uint) ([]domain.User, error)
}

type classService struct {
	dao dao.ClassDao
}

func NewClassService(dao dao.ClassDao) ClassService {
	return &classService{
		dao: dao,
	}
}

func (cs *classService) CreateClass(ctx context.Context, teacherID uint, className string) (domain.Class, error) {
	inviteCode := tool.GenerateRandomCode(8)
	class := domain.Class{
		ClassName:  className,
		TeacherID:  teacherID,
		InviteCode: inviteCode,
	}

	res, err := cs.dao.CreateClass(ctx, class)
	if err != nil {
		return domain.Class{}, errors.CreateClassError(err)
	}
	return res, nil
}

func (cs *classService) JoinClass(ctx context.Context, userID uint, inviteCode string) error {
	err := cs.dao.JoinClass(ctx, inviteCode, userID)
	if err != nil {
		return errors.JoinClassError(err)
	}
	return nil
}

func (cs *classService) QuitClass(ctx context.Context, userID uint, classID uint) error {
	err := cs.dao.QuitClass(ctx, classID, userID)
	if err != nil {
		return errors.QuitClassError(err)
	}
	return nil
}

func (cs *classService) GetStudentClasses(ctx context.Context, userID uint) ([]domain.Class, error) {
	classes, err := cs.dao.GetStudentClasses(ctx, userID)
	if err != nil {
		return nil, errors.GetStudentClassesError(err)
	}
	return classes, nil
}

func (cs *classService) GetTeacherClasses(ctx context.Context, teacherID uint) ([]domain.Class, error) {
	classes, err := cs.dao.GetTeacherClasses(ctx, teacherID)
	if err != nil {
		return nil, errors.GetTeacherClassesError(err)
	}
	return classes, nil
}

func (cs *classService) GetClassByID(ctx context.Context, classID uint) (domain.Class, error) {
	class, err := cs.dao.GetClassByID(ctx, classID)
	if err != nil {
		return domain.Class{}, errors.GetClassInfoError(err)
	}
	return class, nil
}

func (cs *classService) GetClassStudents(ctx context.Context, classID uint) ([]domain.User, error) {
	users, err := cs.dao.GetClassStudents(ctx, classID)
	if err != nil {
		return nil, errors.GetClassStudentsError(err)
	}
	return users, nil
}
