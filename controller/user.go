package controller

import (
	"github.com/Knowpals/Knowpals-be-go/api"
	"github.com/Knowpals/Knowpals-be-go/api/user"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/Knowpals/Knowpals-be-go/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SendCode(c *gin.Context, request user.SendCodeRequest) (api.Response, error)
	Register(c *gin.Context, request user.RegisterRequest) (api.Response, error)
	LoginByPassword(c *gin.Context, request user.LoginByPassword) (api.Response, error)
	LoginByVerifyCode(c *gin.Context, request user.LoginByVerifyCode) (api.Response, error)
	GetUserByID(c *gin.Context, request user.GetUserRequest) (api.Response, error)
}

type userController struct {
	jwt     *ijwt.JwtHandler
	service service.UserService
}

func NewUserController(jwt *ijwt.JwtHandler, service service.UserService) UserController {
	return &userController{
		jwt:     jwt,
		service: service,
	}
}

// SendCode 发送验证码
//
//	@Summary		发送验证码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		user.SendCodeRequest								true	"注册参数"
//	@Router			/api/v1/user/sendCode [post]
func (uc *userController) SendCode(c *gin.Context, request user.SendCodeRequest) (api.Response, error) {
	err := uc.service.SendCode(c, request.Email)
	if err != nil {
		return api.Response{}, err
	}
	return api.SuccessResp(nil), nil
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		user.RegisterRequest								true	"注册参数"
//	@Router			/api/v1/user/register [post]
func (uc *userController) Register(c *gin.Context, request user.RegisterRequest) (api.Response, error) {
	u := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Role:     domain.RoleType(request.Role),
		Email:    request.Email,
	}
	err := uc.service.Register(c.Request.Context(), u, request.Code)
	if err != nil {
		return api.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return api.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.RegisterResp{Token: token}

	return api.SuccessResp(resp), nil
}

// LoginByPassword 用户密码登录
//
//	@Summary		用户密码登录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		user.LoginByPassword								true	"登陆参数"
//	@Router			/api/v1/user/loginByPassword [post]
func (uc *userController) LoginByPassword(c *gin.Context, request user.LoginByPassword) (api.Response, error) {
	u, err := uc.service.LoginByPassword(c, request.Email, request.Password)
	if err != nil {
		return api.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return api.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.LoginResp{Token: token}
	return api.SuccessResp(resp), nil
}

// LoginByVerifyCode 用户验证码登录
//
//	@Summary		用户验证码登录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		user.LoginByVerifyCode								true	"登陆参数"
//	@Router			/api/v1/user/loginByCode [post]
func (uc *userController) LoginByVerifyCode(c *gin.Context, request user.LoginByVerifyCode) (api.Response, error) {
	u, err := uc.service.LoginByVerifyCode(c, request.Email, request.VerifyCode)
	if err != nil {
		return api.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return api.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.LoginResp{Token: token}
	return api.SuccessResp(resp), nil
}

// GetUserByID 根据id查询用户
//
//		@Summary		根据id查询用户
//		@Tags			user
//		@Accept			json
//		@Produce		json
//	 	@Param 			id 		path 		int 		true 		"用户ID"
//		@Router			/api/v1/user/getUser/{id} [get]
//		@Security 		ApiKeyAuth
func (uc *userController) GetUserByID(c *gin.Context, request user.GetUserRequest) (api.Response, error) {
	u, err := uc.service.GetUserByID(c, request.ID)
	if err != nil {
		return api.Response{}, errors.GetUserError(err)
	}

	return api.SuccessResp(u), nil
}
