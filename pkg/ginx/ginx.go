package ginx

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	http2 "github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/gin-gonic/gin"
)

const CTX = "claim"

func SetClaim(c *gin.Context, claim ijwt.UserClaim) {
	c.Set(CTX, claim)
}

func GetClaim(c *gin.Context) (ijwt.UserClaim, error) {
	res, ok := c.Get(CTX)
	if !ok {
		c.Error(errors.New("claim不存在"))
		return ijwt.UserClaim{}, errors.New("claim 不存在")
	}

	claim, ok := res.(ijwt.UserClaim)
	if !ok {
		c.Error(errors.New("claim断言失败"))
		return ijwt.UserClaim{}, errors.New("claim断言失败")
	}
	return claim, nil
}

// ginx主要是为了封装controller
func WrapReq[Req any](fn func(*gin.Context, Req) (http2.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//检查前置中间件是否有错误
		if len(ctx.Errors) > 0 {
			return
		}

		var req Req
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("请求参数错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			ctx.Error(err) //记录错误到ctx，以便后续中间件进行日志处理

			customErr := errorx.ToCustomError(err)

			ctx.JSON(customErr.HttpCode, http2.Response{
				Code:    customErr.Code,
				Message: customErr.Msg,
				Data:    nil,
			})
			return
		}

		//默认的成功响应
		ctx.JSON(ctx.Writer.Status(), res)

	}
}

func WrapUri[Req any](fn func(*gin.Context, Req) (http2.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		err := ctx.BindUri(&req)
		//这里可以直接返回，不需要ctx，因为中间件已经过完了，参数错误也不需要记录日志
		if err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("请求参数错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			//这里需要把错误加到ctx上，因为要后续日志中间件记录日志
			ctx.Error(err)

			//尝试解析成customErr
			e := errorx.ToCustomError(err)
			ctx.JSON(e.HttpCode, http2.Response{
				Code:    e.Code,
				Message: e.Msg,
				Data:    nil,
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func WrapReqAndClaim[Req any](fn func(*gin.Context, Req, ijwt.UserClaim) (http2.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}

		claim, err := GetClaim(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("鉴权错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		var req Req
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("请求参数错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		res, err := fn(ctx, req, claim)
		if err != nil {
			ctx.Error(err)
			customErr := errorx.ToCustomError(err)
			ctx.JSON(customErr.HttpCode, http2.Response{
				Code:    customErr.Code,
				Message: customErr.Msg,
				Data:    nil,
			})
			return
		}
		ctx.JSON(ctx.Writer.Status(), res)
	}
}

func WrapUriAndClaim[Req any](fn func(*gin.Context, Req, ijwt.UserClaim) (http2.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}

		claim, err := GetClaim(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("鉴权错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		var req Req
		err = ctx.BindUri(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("请求参数错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		res, err := fn(ctx, req, claim)
		if err != nil {
			ctx.Error(err)
			e := errorx.ToCustomError(err)
			ctx.JSON(e.HttpCode, http2.Response{
				Code:    e.Code,
				Message: e.Msg,
				Data:    nil,
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func WrapClaim(fn func(*gin.Context, ijwt.UserClaim) (http2.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Errors) > 0 {
			return
		}
		claim, err := GetClaim(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("鉴权错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}

		res, err := fn(ctx, claim)
		if err != nil {
			ctx.Error(err)
			customErr := errorx.ToCustomError(err)
			ctx.JSON(customErr.HttpCode, http2.Response{
				Code:    customErr.Code,
				Message: customErr.Msg,
				Data:    nil,
			})
			return
		}
		ctx.JSON(ctx.Writer.Status(), res)
	}
}

func WrapFormDataAndClaim[Req any](fn func(*gin.Context, Req, multipart.File, *multipart.FileHeader, ijwt.UserClaim) (http2.Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := GetClaim(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("鉴权错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}
		var req Req
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("请求参数错误: %v", err.Error()),
				Data:    nil,
			})
			return
		}
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, http2.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("获取文件失败: %v", err.Error()),
				Data:    nil,
			})
			return
		}
		defer file.Close()

		resp, err := fn(c, req, file, fileHeader, claim)
		if err != nil {
			c.Error(err)
			e := errorx.ToCustomError(err)
			c.JSON(e.HttpCode, http2.Response{
				Code:    e.Code,
				Message: e.Msg,
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, resp)

	}
}
