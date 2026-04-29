package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (cm *CorsMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的请求头
		AllowHeaders: []string{"Content-Type", "Authorization", "Origin"},
		// 是否允许携带凭证（如 Cookies）
		AllowCredentials: true,
		// 解决跨域问题,这个地方允许所有请求跨域了,之后要改成允许前端的请求,比如localhost
		AllowOriginFunc: func(origin string) bool {
			//暂时允许所有跨域请求,根据需要进行调整
			return true
		},

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}
