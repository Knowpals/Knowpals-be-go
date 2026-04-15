package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/statistic"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterStatisticRoute(r *gin.RouterGroup, sc statistic.StatController, auth *middleware.AuthMiddleware) {
	c := r.Group("/stat")
	c.Use(auth.MiddlewareFunc())
	{
		c.GET("/student/overview", ginx.WrapClaim(sc.GetStudentOverview))
		c.GET("/student/:video_id", ginx.WrapUri(sc.GetStudentStat))
		c.GET("/class", ginx.WrapReq(sc.GetClassStat))
	}
}

