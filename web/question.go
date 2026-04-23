package web

import (
	"github.com/Knowpals/Knowpals-be-go/controller/question"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

func RegisterQuestionRoute(r *gin.RouterGroup, qc question.QuestionController, auth *middleware.AuthMiddleware) {
	c := r.Group("/question")
	c.Use(auth.MiddlewareFunc())
	{
		c.POST("/answer", ginx.WrapReq(qc.AnswerQuestion))
		c.GET("/generate/:video_id", ginx.WrapUri(qc.GenerateVideoExercise))

		// teacher review
		c.GET("/review/:video_id", ginx.WrapUri(qc.ReviewListQuestions))
		c.POST("/review/add", ginx.WrapReq(qc.ReviewAddQuestion))
		// update uses uri + json body, bind json inside controller
		c.PUT("/review/:question_id", ginx.WrapUri(qc.ReviewUpdateQuestion))
		c.DELETE("/review/:question_id", ginx.WrapUri(qc.ReviewDeleteQuestion))
	}
}

