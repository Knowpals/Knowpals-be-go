package user

import (
	"context"
	errors1 "errors"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/infra/email"
	"github.com/Knowpals/Knowpals-be-go/pkg/otelx"
	"github.com/Knowpals/Knowpals-be-go/repository/cache"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/tool"
)

type UserService interface {
	SendCode(ctx context.Context, targetAddr string) error
	Register(ctx context.Context, user *domain.User, code string) error
	//密码登陆
	LoginByPassword(ctx context.Context, email string, password string) (domain.User, error)
	LoginByVerifyCode(ctx context.Context, email string, code string) (domain.User, error)
	GetUserByID(ctx context.Context, id uint) (domain.User, error)
}

type userService struct {
	dao   dao.UserDao
	cache cache.AuthCache
	email *email.EmailClient
}

func NewUserService(dao dao.UserDao, email *email.EmailClient, cache cache.AuthCache) UserService {
	return &userService{
		dao:   dao,
		email: email,
		cache: cache,
	}
}

func (us *userService) GetUserByID(ctx context.Context, id uint) (domain.User, error) {
	user, err := us.dao.GetUserByID(ctx, id)
	if err != nil {
		otelx.RecordError(ctx, err)
		return domain.User{}, err
	}
	return user, nil
}

func (us *userService) SendCode(ctx context.Context, targetAddr string) error {
	//生成验证码
	code := tool.GenerateVerifyCode()

	err := us.cache.SetCode(ctx, code, targetAddr)
	if err != nil {
		return errors.SendVerifyCodeError(err)
	}
	//发送邮件
	err = us.email.SendEmail(targetAddr, code)
	if err != nil {
		return errors.SendVerifyCodeError(err)
	}
	return nil
}

// TODO：验证码不让连续触发
func (us *userService) Register(ctx context.Context, user *domain.User, code string) error {
	//校验
	c, err := us.cache.GetCode(ctx, user.Email)
	if err != nil {
		return errors.VerifyCodeError(err)
	}
	if c != code {
		return errors.VerifyCodeError(errors1.New("验证码校验错误"))
	}
	//写入数据库，注册成功
	err = us.dao.AddUser(ctx, user)
	if err != nil {
		return errors.RegisterUserError(err)
	}
	return nil
}

func (us *userService) LoginByPassword(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := us.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, errors.LoginError(err)
	}

	if u.Password != password {
		return domain.User{}, errors.LoginError(errors1.New("密码错误"))
	}

	return u, nil

}

func (us *userService) LoginByVerifyCode(ctx context.Context, email string, code string) (domain.User, error) {
	//验证是否注册过，存在此用户
	u, err := us.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, errors.LoginError(err)
	}

	c, err := us.cache.GetCode(ctx, email)
	if err != nil {
		return domain.User{}, errors.LoginError(err)
	}

	if c != code {
		return domain.User{}, errors.LoginError(errors1.New("验证码错误"))
	}

	return u, nil

}
