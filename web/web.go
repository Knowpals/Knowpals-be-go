package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/class"
	"github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/controller/video"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/gin-gonic/gin"
)

func NewGinEngine(
	uc user.UserController,
	cc class.ClassController,
	vc video.VideoController,
	auth *middleware.AuthMiddleware,
	log *middleware.LoggerMiddleware,
	// otel *middleware.OtelMiddleware,
) *gin.Engine {
	r := gin.Default()
	r.Use(
		//otel.MiddlewareFunc(),
		log.MiddlewareFunc(),
	)
	apiV1 := r.Group("/api/v1")
	RegisterUserRoute(apiV1, uc, auth)
	RegisterClassRoute(apiV1, cc, auth)
	RegisterVideoRoute(apiV1, vc, auth)
	return r
}
