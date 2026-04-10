package user

import (
	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/user"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	user2 "github.com/Knowpals/Knowpals-be-go/service/user"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	// SendCode 发送验证码
	SendCode(c *gin.Context, request user.SendCodeRequest) (http.Response, error)
	// Register 用户注册
	Register(c *gin.Context, request user.RegisterRequest) (http.Response, error)
	// LoginByPassword 密码登录
	LoginByPassword(c *gin.Context, request user.LoginByPassword) (http.Response, error)
	// LoginByVerifyCode 验证码登录
	LoginByVerifyCode(c *gin.Context, request user.LoginByVerifyCode) (http.Response, error)
	// GetUserByID 通过id查找用户
	GetUserByID(c *gin.Context, request user.GetUserRequest) (http.Response, error)
}

type userController struct {
	jwt     *ijwt.JwtHandler
	service user2.UserService
}

func NewUserController(jwt *ijwt.JwtHandler, service user2.UserService) UserController {
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
func (uc *userController) SendCode(c *gin.Context, request user.SendCodeRequest) (http.Response, error) {
	err := uc.service.SendCode(c, request.Email)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// Register 用户注册
//
//		@Summary		用户注册
//		@Tags			user
//		@Accept			json
//		@Produce		json
//		@Param			request	body		user.RegisterRequest								true	"注册参数"
//	 @Success 200 {object} http.Response{data=user.RegisterResp} "成功"
//		@Router			/api/v1/user/register [post]
func (uc *userController) Register(c *gin.Context, request user.RegisterRequest) (http.Response, error) {
	u := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Role:     domain.RoleType(request.Role),
		Email:    request.Email,
	}
	err := uc.service.Register(c.Request.Context(), u, request.Code)
	if err != nil {
		return http.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return http.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.RegisterResp{Token: token}

	return http.Success(resp), nil
}

// LoginByPassword 用户密码登录
//
//		@Summary		用户密码登录
//		@Tags			user
//		@Accept			json
//		@Produce		json
//		@Param			request	body		user.LoginByPassword								true	"登陆参数"
//	 @Success 200 {object} http.Response{data=user.LoginResp} "成功"
//		@Router			/api/v1/user/loginByPassword [post]
func (uc *userController) LoginByPassword(c *gin.Context, request user.LoginByPassword) (http.Response, error) {
	u, err := uc.service.LoginByPassword(c, request.Email, request.Password)
	if err != nil {
		return http.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return http.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.LoginResp{Token: token}
	return http.Success(resp), nil
}

// LoginByVerifyCode 用户验证码登录
//
//		@Summary		用户验证码登录
//		@Tags			user
//		@Accept			json
//		@Produce		json
//		@Param			request	body		user.LoginByVerifyCode								true	"登陆参数"
//	 	@Success 200 {object} http.Response{data=user.LoginResp} "成功"
//		@Router			/api/v1/user/loginByCode [post]
func (uc *userController) LoginByVerifyCode(c *gin.Context, request user.LoginByVerifyCode) (http.Response, error) {
	u, err := uc.service.LoginByVerifyCode(c, request.Email, request.VerifyCode)
	if err != nil {
		return http.Response{}, err
	}

	token, err := uc.jwt.GenerateToken(u.ID, u.Username, u.Password, u.Email, string(u.Role))
	if err != nil {
		return http.Response{}, errors.TokenGenerateError(err)
	}

	resp := user.LoginResp{Token: token}
	return http.Success(resp), nil
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
func (uc *userController) GetUserByID(c *gin.Context, request user.GetUserRequest) (http.Response, error) {
	u, err := uc.service.GetUserByID(c, request.ID)
	if err != nil {
		return http.Response{}, errors.GetUserError(err)
	}

	return http.Success(u), nil
}
