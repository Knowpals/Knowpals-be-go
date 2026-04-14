package question

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionController interface {
	// AnswerQuestion 学生回答问题
	AnswerQuestion(c *gin.Context, req question.AnswerQuestionReq) (http.Response, error)
	// GenerateVideoExercise 生成视频任务的课后习题
	GenerateVideoExercise(c *gin.Context, req question.GenerateVideoExerciseReq) (http.Response, error)
}

type questionController struct {
	db *gorm.DB
}

func NewQuestionController(db *gorm.DB) QuestionController {
	return &questionController{db: db}
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
	if len(req.StudentAnswers) == 0 {
		return http.Success(question.AnswerQuestionResp{Results: []question.Result{}}), nil
	}

	qids := make([]uint, 0, len(req.StudentAnswers))
	for _, a := range req.StudentAnswers {
		qids = append(qids, a.QuestionID)
	}

	var questions []model.Question
	if err := qc.db.WithContext(c).
		Where("id IN ? AND video_id = ?", qids, req.VideoID).
		Find(&questions).Error; err != nil {
		return http.Response{}, err
	}
	qmap := make(map[uint]model.Question, len(questions))
	for _, q := range questions {
		qmap[q.ID] = q
	}

	results := make([]question.Result, 0, len(req.StudentAnswers))
	answersToInsert := make([]model.StudentAnswer, 0, len(req.StudentAnswers))

	normalize := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

	for _, a := range req.StudentAnswers {
		q, ok := qmap[a.QuestionID]
		if !ok {
			continue
		}
		isCorrect := normalize(a.Answer) == normalize(q.Answer)
		results = append(results, question.Result{
			QuestionID: q.ID,
			IsCorrect:  isCorrect,
			Answer:     q.Answer,
			Analysis:   q.Analysis,
		})
		answersToInsert = append(answersToInsert, model.StudentAnswer{
			StudentID:  claim.ID,
			QuestionID: q.ID,
			IsCorrect:  isCorrect,
			TimeCost:   float64(a.TimeCost),
			VideoID:    req.VideoID,
			SegmentID:  q.SegmentID,
		})
	}

	if err := qc.db.WithContext(c).CreateInBatches(&answersToInsert, 200).Error; err != nil {
		return http.Response{}, err
	}

	return http.Success(question.AnswerQuestionResp{Results: results}), nil
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
	// 已由 pipeline 生成并落库，此接口可作为查询/重跑入口；先返回当前题目列表
	var qs []model.Question
	if err := qc.db.WithContext(c).Where("video_id = ?", req.VideoID).Find(&qs).Error; err != nil {
		return http.Response{}, err
	}
	out := make([]question.Question, 0, len(qs))
	for _, q := range qs {
		opts := []string{}
		if q.Options != nil && *q.Options != "" {
			_ = json.Unmarshal([]byte(*q.Options), &opts)
		}
		out = append(out, question.Question{
			ID:        q.ID,
			Type:      q.Type,
			Content:   q.Content,
			Options:   opts,
			Answer:    q.Answer,
			Analysis:  q.Analysis,
			SegmentID: q.SegmentID,
		})
	}
	return http.Success(question.GenerateVideoExerciseResp{Questions: out}), nil
}
