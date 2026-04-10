package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type ClassDao interface {
	CreateClass(ctx context.Context, class domain.Class) (domain.Class, error)
	JoinClass(ctx context.Context, inviteCode string, studentID uint) error
	QuitClass(ctx context.Context, classID uint, studentID uint) error
	GetStudentClasses(ctx context.Context, studentID uint) ([]domain.Class, error)
	GetTeacherClasses(ctx context.Context, teacherID uint) ([]domain.Class, error)
	GetClassByID(ctx context.Context, classID uint) (domain.Class, error)
	GetClassStudents(ctx context.Context, classID uint) ([]domain.User, error)
}

type classDao struct {
	db *gorm.DB
}

func NewClassDao(db *gorm.DB) ClassDao {
	return &classDao{db: db}
}

func (cd *classDao) CreateClass(ctx context.Context, class domain.Class) (domain.Class, error) {
	classModel := &model.Class{
		ClassName:  class.ClassName,
		TeacherID:  class.TeacherID,
		InviteCode: class.InviteCode,
	}
	err := cd.db.WithContext(ctx).Create(classModel).Error
	if err != nil {
		return domain.Class{}, err
	}

	return domain.Class{
		ID:         classModel.ID,
		ClassName:  classModel.ClassName,
		TeacherID:  classModel.TeacherID,
		InviteCode: classModel.InviteCode,
	}, err
}

func (cd *classDao) JoinClass(ctx context.Context, inviteCode string, studentID uint) error {
	//先根据inviteCode查询班级id，再插入记录
	return cd.db.WithContext(ctx).Exec(`
		insert into class_students (class_id,student_id)
		select id,?
		from classes
		where invite_code=?
	`, studentID, inviteCode).Error
}

func (cd *classDao) QuitClass(ctx context.Context, classID uint, studentID uint) error {
	return cd.db.WithContext(ctx).
		Table("class_students").
		Where("class_id=? AND student_id=?", classID, studentID).
		Delete(&model.ClassStudent{}).
		Error
}

func (cd *classDao) GetStudentClasses(ctx context.Context, studentID uint) ([]domain.Class, error) {
	//查学生加入的所有班级
	var classListModels []model.Class
	err := cd.db.WithContext(ctx).
		Table("classes").
		Joins("join class_students on classes.id=class_students.class_id").
		Where("class_students.student_id=?", studentID).
		Find(&classListModels).Error
	if err != nil {
		return nil, err
	}
	classList := make([]domain.Class, len(classListModels))
	for index, m := range classListModels {
		classList[index] = domain.Class{
			ID:         m.ID,
			ClassName:  m.ClassName,
			TeacherID:  m.TeacherID,
			InviteCode: m.InviteCode,
		}
	}
	return classList, nil
}

func (cd *classDao) GetTeacherClasses(ctx context.Context, teacherID uint) ([]domain.Class, error) {
	//查老师创建的所有班级
	var classListModels []model.Class
	err := cd.db.WithContext(ctx).
		Table("classes").
		Where("teacher_id=?", teacherID).
		Find(&classListModels).
		Error
	if err != nil {
		return nil, err
	}

	classList := make([]domain.Class, len(classListModels))
	for i, m := range classListModels {
		classList[i] = domain.Class{
			ID:         m.ID,
			TeacherID:  teacherID,
			InviteCode: m.InviteCode,
			ClassName:  m.ClassName,
		}
	}
	return classList, nil
}

func (cd *classDao) GetClassByID(ctx context.Context, classID uint) (domain.Class, error) {
	var classModel model.Class
	err := cd.db.WithContext(ctx).
		Table("classes").
		Where("id=?", classID).
		First(&classModel).
		Error
	if err != nil {
		return domain.Class{}, err
	}
	return domain.Class{
		ID:         classModel.ID,
		ClassName:  classModel.ClassName,
		TeacherID:  classModel.TeacherID,
		InviteCode: classModel.InviteCode,
	}, nil
}

func (cd *classDao) GetClassStudents(ctx context.Context, classID uint) ([]domain.User, error) {
	var userModels []model.User
	err := cd.db.WithContext(ctx).
		Table("users").
		Joins("join class_students on class_students.student_id=users.id").
		Where("class_students.class_id=?", classID).
		Find(&userModels).
		Error

	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(userModels))
	for i, m := range userModels {
		users[i] = domain.User{
			ID:       m.ID,
			Username: m.Username,
			Email:    m.Email,
			Password: m.Password,
			Role:     domain.RoleType(m.Role),
		}
	}
	return users, nil
}
