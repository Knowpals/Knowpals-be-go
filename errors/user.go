package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	//user
	RegisterUserErrorCode = 50000 + iota
	TokenGenerateErrorCode
	GetUserErrorCode
	SendVerifyCodeErrorCode
	VerifyCodeErrorCode
	LoginErrorCode
)

const (
	RegisterRequestErrorCode = 40000 + iota
)

var (
	//user
	RegisterUserError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, RegisterUserErrorCode, "用户注册失败", err)
	}

	RegisterRequestError = func() error {
		return errorx.New(http.StatusBadRequest, RegisterRequestErrorCode, "用户注册参数错误", nil)
	}

	TokenGenerateError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, TokenGenerateErrorCode, "token生成失败", err)
	}

	GetUserError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetUserErrorCode, "获取学生失败", err)
	}

	SendVerifyCodeError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, SendVerifyCodeErrorCode, "发送验证码失败", err)
	}

	VerifyCodeError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, VerifyCodeErrorCode, "验证码错误", err)
	}

	LoginError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, LoginErrorCode, "登陆失败", err)
	}
)
