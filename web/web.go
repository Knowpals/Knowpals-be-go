package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/behavior"
	"github.com/Knowpals/Knowpals-be-go/controller/class"
	"github.com/Knowpals/Knowpals-be-go/controller/question"
	"github.com/Knowpals/Knowpals-be-go/controller/statistic"
	"github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/controller/video"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/gin-gonic/gin"
)

func NewGinEngine(
	uc user.UserController,
	cc class.ClassController,
	vc video.VideoController,
	qc question.QuestionController,
	bc behavior.BehaviorController,
	sc statistic.StatController,
	auth *middleware.AuthMiddleware,
	log *middleware.LoggerMiddleware,
	cors *middleware.CorsMiddleware,
) *gin.Engine {
	r := gin.Default()
	r.Use(
		log.MiddlewareFunc(),
		cors.MiddlewareFunc(),
	)
	apiV1 := r.Group("/api/v1")
	RegisterUserRoute(apiV1, uc, auth)
	RegisterClassRoute(apiV1, cc, auth)
	RegisterVideoRoute(apiV1, vc, auth)
	RegisterQuestionRoute(apiV1, qc, auth)
	RegisterBehaviorRoute(apiV1, bc, auth)
	RegisterStatisticRoute(apiV1, sc, auth)
	return r
}
