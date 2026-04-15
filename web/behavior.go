package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/behavior"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterBehaviorRoute(r *gin.RouterGroup, bc behavior.BehaviorController, auth *middleware.AuthMiddleware) {
	c := r.Group("/behavior")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/record", ginx.WrapReqAndClaim(bc.RecordAction))
		c.POST("/update-progress", ginx.WrapReqAndClaim(bc.UpdateVideoProgress))
		c.GET("/class-progress/:class_id/:status", ginx.WrapUriAndClaim(bc.GetClassVideoProgress))
	}
}
