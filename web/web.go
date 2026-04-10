package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/gin-gonic/gin"
)

func NewGinEngine(
	uc user.UserController,
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

	return r
}
