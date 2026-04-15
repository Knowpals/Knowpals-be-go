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
		c.GET("/getDetail/:video_id", ginx.WrapUri(controller.GetVideoDetail))
		c.POST("/post-to-class", ginx.WrapReq(controller.PostVideoToClass))
		c.GET("/getTasks/:class_id", ginx.WrapUri(controller.GetClassVideoTasks))
		c.POST("/task/process", ginx.WrapReq(controller.GetTaskUploadingProcess))
		c.GET("/my-uploaded", ginx.WrapClaim(controller.GetMyUploadedVideos))
	}

}
