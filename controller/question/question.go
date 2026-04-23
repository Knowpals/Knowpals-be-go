package question

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	questionSvc "github.com/Knowpals/Knowpals-be-go/service/question"
	"github.com/gin-gonic/gin"
)

type QuestionController interface {
	// AnswerQuestion 学生回答问题
	AnswerQuestion(c *gin.Context, req question.AnswerQuestionReq) (http.Response, error)
	// GenerateVideoExercise 生成视频任务的课后习题
	GenerateVideoExercise(c *gin.Context, req question.GenerateVideoExerciseReq) (http.Response, error)
	// ReviewListQuestions 老师查看待审核题目
	ReviewListQuestions(c *gin.Context, req question.ReviewListReq) (http.Response, error)
	// ReviewAddQuestion 老师按时间点添加题目
	ReviewAddQuestion(c *gin.Context, req question.ReviewAddReq) (http.Response, error)
	// ReviewUpdateQuestion 老师修改题目
	ReviewUpdateQuestion(c *gin.Context, req question.ReviewUpdateReq) (http.Response, error)
	// ReviewDeleteQuestion 老师删除题目
	ReviewDeleteQuestion(c *gin.Context, req question.ReviewDeleteReq) (http.Response, error)
}

type questionController struct {
	svc questionSvc.Service
}

func NewQuestionController(svc questionSvc.Service) QuestionController {
	return &questionController{svc: svc}
}

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
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := qc.svc.AnswerQuestion(c, claim.ID, req)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
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
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	resp, err := qc.svc.GenerateVideoExercise(c, claim.ID, domain.RoleType(claim.Role), req.VideoID)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}
