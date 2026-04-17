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
	ForgotPassword(ctx context.Context, email string, code string, newPassword string) error
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
		return domain.User{}, errors.GetUserError(err)
	}
	return user, nil
}

func (us *userService) SendCode(ctx context.Context, targetAddr string) error {
	//生成验证码
	code := tool.GenerateRandomCode(6)

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
	_, isExist, err := us.dao.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return errors.RegisterUserError(err)
	}
	if isExist {
		return errors.RegisterUserError(errors1.New("已有账号，不可重复注册"))
	}
	//写入数据库，注册成功
	err = us.dao.AddUser(ctx, user)
	if err != nil {
		return errors.RegisterUserError(err)
	}
	return nil
}

func (us *userService) LoginByPassword(ctx context.Context, email string, password string) (domain.User, error) {
	u, isExist, err := us.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, errors.LoginError(err)
	}

	if !isExist {
		return domain.User{}, errors.LoginError(errors1.New("用户不存在"))
	}

	if u.Password != password {
		return domain.User{}, errors.LoginError(errors1.New("密码错误"))
	}

	return u, nil
}

func (us *userService) LoginByVerifyCode(ctx context.Context, email string, code string) (domain.User, error) {
	//验证是否注册过，存在此用户
	u, isExist, err := us.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, errors.LoginError(err)
	}
	if !isExist {
		return domain.User{}, errors.LoginError(errors1.New("用户不存在"))
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

func (us *userService) ForgotPassword(ctx context.Context, email string, code string, newPassword string) error {
	c, err := us.cache.GetCode(ctx, email)
	if err != nil {
		return errors.ForgotPasswordError(err)
	}
	if c != code {
		return errors.ForgotPasswordError(errors1.New("验证码错误"))
	}
	// 确认用户存在
	_, isExist, err := us.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.ForgotPasswordError(err)
	}
	if !isExist {
		return errors.ForgotPasswordError(errors1.New("用户不存在"))
	}
	if err := us.dao.UpdatePasswordByEmail(ctx, email, newPassword); err != nil {
		return errors.ForgotPasswordError(err)
	}
	return nil
}
