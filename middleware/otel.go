package middleware

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type OtelMiddleware struct {
	serviceName string
}

func NewOtelMiddleware(conf *config.Config) *OtelMiddleware {
	return &OtelMiddleware{
		serviceName: conf.Otel.ServiceName,
	}
}

func (om *OtelMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return otelgin.Middleware(om.serviceName)
}
