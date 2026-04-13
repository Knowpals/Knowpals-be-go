package web

import (
	class2 "github.com/Knowpals/Knowpals-be-go/controller/class"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterClassRoute(r *gin.RouterGroup, classController class2.ClassController, auth *middleware.AuthMiddleware) {
	c := r.Group("/class")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/create", ginx.WrapReqAndClaim(classController.CreateClass))
		c.POST("/join", ginx.WrapReqAndClaim(classController.JoinClass))
		c.POST("/quit/:class_id", ginx.WrapUriAndClaim(classController.QuitClass))
		c.GET("/info/:class_id", ginx.WrapUriAndClaim(classController.GetClassInfo))
		c.GET("/my-created", ginx.WrapClaim(classController.GetMyCreatedClasses))
		c.GET("/my-joined", ginx.WrapClaim(classController.GetMyJoinedClasses))
		c.GET("/students/:class_id", ginx.WrapUri(classController.GetClassStudents))
	}
}
