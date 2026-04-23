package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/agent"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterAgentRoute(r *gin.RouterGroup, ac agent.AgentController, auth *middleware.AuthMiddleware) {
	c := r.Group("/agent")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/chat", ginx.WrapReq(ac.Chat))
		c.GET("/history", ginx.WrapReq(ac.GetChatHistory))
		c.POST("/quiz", ginx.WrapReq(ac.GenerateQuiz))
		c.POST("/report", ginx.WrapReq(ac.GenerateReport))
		c.GET("/report", ginx.WrapReq(ac.GetReport))
	}
}

