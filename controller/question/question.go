package question

import (
	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/gin-gonic/gin"
)

type QuestionController interface {
	// AnswerQuestion 学生回答问题
	AnswerQuestion(c *gin.Context, req question.AnswerQuestionReq) (http.Response, error)
	// GenerateVideoExercise 生成视频任务的课后习题
	GenerateVideoExercise(c *gin.Context, req question.GenerateVideoExerciseReq) (http.Response, error)
}

type questionController struct{}

// AnswerQuestion 学生提交答案
// @Summary 学生批量提交答题答案
// @Description 学生提交单个视频/片段内的题目答案，系统自动批改并返回结果
// @Tags question
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body question.AnswerQuestionReq true "学生答题参数"
// @Success 200 {object} http.Response{data=question.AnswerQuestionResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/question/answer [post]
func (qc *questionController) AnswerQuestion(c *gin.Context, req question.AnswerQuestionReq) (http.Response, error) {
	return http.Response{}, nil
}

// GenerateVideoExercise 生成视频课后习题
// @Summary AI生成视频课后习题
// @Description 根据视频ID自动生成课后练习题，供教师确认后发布
// @Tags question
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id path uint true "视频ID参数"
// @Success 200 {object} http.Response{data=question.GenerateVideoExerciseResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/question/generate/{video_id} [get]
func (qc *questionController) GenerateVideoExercise(c *gin.Context, req question.GenerateVideoExerciseReq) (http.Response, error) {
	return http.Response{}, nil
}
