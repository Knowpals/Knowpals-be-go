package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/review"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterReviewRoute(r *gin.RouterGroup, rc review.ReviewController, auth *middleware.AuthMiddleware) {
	c := r.Group("/video")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/:video_id/review/start", ginx.WrapUri(rc.Start))
		c.POST("/:video_id/review/publish", ginx.WrapUri(rc.Publish))
	}
}

