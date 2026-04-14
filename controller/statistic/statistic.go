package statistic

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/statistic"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	statService "github.com/Knowpals/Knowpals-be-go/service/statistic"
	"github.com/gin-gonic/gin"
)

type StatController interface {
	// GetStudentStat 获取学生的学习情况
	GetStudentStat(c *gin.Context, req statistic.GetStudentStatReq) (http.Response, error)
	// GetClassStat 获取班级学习情况
	GetClassStat(c *gin.Context, req statistic.GetClassStatReq) (http.Response, error)
}

type statController struct {
	svc statService.StatService
}

func NewStatController(svc statService.StatService) StatController {
	return &statController{svc: svc}
}

// GetStudentStat 获取学生个人学习统计
// @Summary 获取学生个人学情统计
// @Description 获取学生单个视频的答题正确率、观看时长、暂停次数、薄弱知识点
// @Tags statistic
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id path int true "视频ID"
// @Success 200 {object} http.Response{data=statistic.GetStudentStatResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/stat/student/{video_id} [get]
func (sc *statController) GetStudentStat(c *gin.Context, req statistic.GetStudentStatReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	statDomain, err := sc.svc.GetStudentStat(c, claim.ID, req.VideoID)
	if err != nil {
		return http.Response{}, err
	}
	weak := make([]statistic.KnowledgePoint, 0, len(statDomain.WeakKnowledgePoints))
	for _, kp := range statDomain.WeakKnowledgePoints {
		weak = append(weak, statistic.KnowledgePoint{
			KnowledgeID: kp.KnowledgeID,
			Title:       kp.Title,
			MasterScore: kp.MasterScore,
		})
	}
	return http.Success(statistic.GetStudentStatResp{
		Status:              statDomain.Status,
		CorrectRate:         statDomain.CorrectRate,
		TimeCost:            statDomain.TimeCost,
		PauseCount:          statDomain.PauseCount,
		WeakKnowledgePoints: weak,
	}), nil
}

// GetClassStat 获取班级学习统计
// @Summary 获取班级整体学情统计
// @Description 获取班级在某个视频下的整体数据：正确率、完成率、薄弱知识点、高频暂停/回放片段
// @Tags statistic
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body statistic.GetClassStatReq true "获取班级情况参数"
// @Success 200 {object} http.Response{data=statistic.GetClassStatResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/stat/class [get]
func (sc *statController) GetClassStat(c *gin.Context, req statistic.GetClassStatReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if !domain.RoleType(claim.Role).IsValid() {
		return http.Response{}, errors.New("无权限")
	}
	statDomain, err := sc.svc.GetClassStat(c, req.ClassID, req.VideoID)
	if err != nil {
		return http.Response{}, err
	}
	topQs := make([]statistic.QuestionStat, 0, len(statDomain.TopQuestions))
	for _, q := range statDomain.TopQuestions {
		topQs = append(topQs, statistic.QuestionStat{
			QuestionID: q.QuestionID,
			Content:    q.Content,
			ErrorRate:  q.ErrorRate,
		})
	}
	pauseActs := make([]statistic.PauseAction, 0, len(statDomain.TopPauseAction))
	for _, a := range statDomain.TopPauseAction {
		pauseActs = append(pauseActs, statistic.PauseAction{
			SegmentID:  a.SegmentID,
			Start:      a.Start,
			End:        a.End,
			PauseCount: a.PauseCount,
		})
	}
	replayActs := make([]statistic.ReplayAction, 0, len(statDomain.TopReplayAction))
	for _, a := range statDomain.TopReplayAction {
		replayActs = append(replayActs, statistic.ReplayAction{
			SegmentID:   a.SegmentID,
			Start:       a.Start,
			End:         a.End,
			ReplayCount: a.ReplayCount,
		})
	}

	return http.Success(statistic.GetClassStatResp{
		Overview: statistic.Overview{
			AverageCorrectRate: statDomain.Overview.AverageCorrectRate,
			AverageTimeCost:    statDomain.Overview.AverageTimeCost,
			CompleteRate:       statDomain.Overview.CompleteRate,
			TotalPauseCount:    statDomain.Overview.TotalPauseCount,
		},
		TopQuestions:    topQs,
		TopPauseAction:  pauseActs,
		TopReplayAction: replayActs,
	}), nil
}
