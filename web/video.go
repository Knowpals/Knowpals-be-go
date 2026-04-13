package web

import (
	video1 "github.com/Knowpals/Knowpals-be-go/controller/video"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterVideoRoute(r *gin.RouterGroup, controller video1.VideoController, auth *middleware.AuthMiddleware) {
	c := r.Group("/video")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/upload", ginx.WrapFormDataAndClaim(controller.UploadVideo))
	}

}
