package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type UserDao interface {
	AddUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id uint) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdatePasswordByEmail(ctx context.Context, email string, newPassword string) error
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{
		db: db,
	}
}

func (ud *userDao) AddUser(ctx context.Context, user *domain.User) error {
	//ctx, span := otelx.Tracer("dao.user").Start(ctx, "UserDao.AddUser")
	//defer span.End()
	//otelx.SetAttributes(ctx, attribute.String("db.operation", "insert"))

	userModel := model.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Role:     string(user.Role),
	}
	err := ud.db.WithContext(ctx).Create(&userModel).Error
	if err != nil {
		return err
	}
	user.ID = userModel.ID
	//otelx.SetAttributes(ctx, attribute.Int64("user.id", int64(user.ID)))
	return nil
}

func (ud *userDao) GetUserByID(ctx context.Context, id uint) (domain.User, error) {
	//ctx, span := otelx.Tracer("dao.user").Start(ctx, "UserDao.GetUserByID")
	//defer span.End()
	//otelx.SetAttributes(ctx,
	//	attribute.String("db.operation", "select"),
	//	attribute.Int64("user.id", int64(id)),
	//)

	var userModel model.User
	err := ud.db.WithContext(ctx).First(&userModel, id).Error
	if err != nil {
		//otelx.RecordError(ctx, err)
		return domain.User{}, err
	}
	role := domain.RoleType(userModel.Role)
	if !role.IsValid() {
		err = errors.New(fmt.Sprintf("角色类型不合法：%s", role))
		//otelx.RecordError(ctx, err)
		return domain.User{}, err
	}
	return domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Password: userModel.Password,
		Email:    userModel.Email,
		Role:     role,
	}, nil
}

func (ud *userDao) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var userModel model.User
	err := ud.db.WithContext(ctx).Where("email=?", email).First(&userModel).Error
	if err != nil {
		return domain.User{}, err
	}
	role := domain.RoleType(userModel.Role)
	if !role.IsValid() {
		err = errors.New(fmt.Sprintf("角色类型不合法：%s", role))
		return domain.User{}, err
	}
	return domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
		Role:     role,
	}, nil
}

func (ud *userDao) DeleteUser(ctx context.Context, id string) error {
	return ud.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (ud *userDao) UpdatePasswordByEmail(ctx context.Context, email string, newPassword string) error {
	return ud.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("password", newPassword).
		Error
}
