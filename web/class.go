package web

import (
	class2 "github.com/Knowpals/Knowpals-be-go/controller/class"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterClassRoute(r *gin.RouterGroup, classController class2.ClassController, auth *middleware.AuthMiddleware) {
	c := r.Group("/class")
	{
		c.POST("/create", auth.MiddlewareFunc(), ginx.WrapReqAndClaim(classController.CreateClass))
		c.POST("/join", auth.MiddlewareFunc(), ginx.WrapReqAndClaim(classController.JoinClass))
		c.POST("/quit/:class_id", auth.MiddlewareFunc(), ginx.WrapUriAndClaim(classController.QuitClass))
		c.GET("/info/:class_id", auth.MiddlewareFunc(), ginx.WrapUriAndClaim(classController.GetClassInfo))
		c.GET("/my-created", auth.MiddlewareFunc(), ginx.WrapClaim(classController.GetMyCreatedClasses))
		c.GET("/my-joined", auth.MiddlewareFunc(), ginx.WrapClaim(classController.GetMyJoinedClasses))
		c.GET("/students/:class_id", auth.MiddlewareFunc(), ginx.WrapUri(classController.GetClassStudents))
	}
}
