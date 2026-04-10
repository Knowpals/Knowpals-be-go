package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

// 这里需要ginx包装
func RegisterUserRoute(r *gin.RouterGroup, userController user.UserController, auth *middleware.AuthMiddleware) {
	c := r.Group("/user")
	{
		c.POST("/register", ginx.WrapReq(userController.Register))
		c.POST("/sendCode", ginx.WrapReq(userController.SendCode))
		c.POST("/loginByPassword", ginx.WrapReq(userController.LoginByPassword))
		c.POST("/loginByCode", ginx.WrapReq(userController.LoginByVerifyCode))
		c.GET("/getUser/:id", auth.MiddlewareFunc(), ginx.WrapUri(userController.GetUserByID))
	}
}
