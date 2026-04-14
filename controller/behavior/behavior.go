package behavior

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/behavior"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	behavior2 "github.com/Knowpals/Knowpals-be-go/service/behavior"
	"github.com/gin-gonic/gin"
)

type BehaviorController interface {
	// RecordAction 记录学生观看行为
	RecordAction(c *gin.Context, req behavior.RecordActionReq, claim ijwt.UserClaim) (http.Response, error)
	// UpdateVideoProgress 更新学生观看进度
	UpdateVideoProgress(c *gin.Context, req behavior.UpdateProgressReq, claim ijwt.UserClaim) (http.Response, error)
	// GetClassVideoProgress 获取学生班级内所有视频观看进度
	GetClassVideoProgress(c *gin.Context, req behavior.GetClassVideoProgressReq, claim ijwt.UserClaim) (http.Response, error)
}

type behaviorController struct {
	svc behavior2.BehaviorService
}

func NewBehaviorController(svc behavior2.BehaviorService) BehaviorController {
	return &behaviorController{svc: svc}
}

// RecordAction 记录播放行为（pause/replay/play）
// @Summary 记录学生视频观看行为
// @Description 记录暂停、回放、播放等行为及时长
// @Tags behavior
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body behavior.RecordActionReq true "行为上报参数"
// @Success 200 {object} http.Response "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/behavior/record [post]
func (bc *behaviorController) RecordAction(c *gin.Context, req behavior.RecordActionReq, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	if err := bc.svc.RecordAction(c, claim.ID, domain.WatchAction{
		VideoID:   req.VideoID,
		SegmentID: req.SegmentID,
		Event:     req.Event,
		Duration:  req.Duration,
	}); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// UpdateVideoProgress 更新视频播放进度
// @Summary 更新学生视频观看进度
// @Description 上报当前播放秒数，更新最大进度与完成状态
// @Tags behavior
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body behavior.UpdateProgressReq true "进度参数"
// @Success 200 {object} http.Response "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/behavior/update-progress [post]
func (bc *behaviorController) UpdateVideoProgress(c *gin.Context, req behavior.UpdateProgressReq, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	if err := bc.svc.UpdateProgress(c, claim.ID, domain.WatchProgress{
		VideoID:    req.VideoID,
		CurrentSec: req.CurrentSec,
	}); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// GetClassVideoProgress 获取班级视频学习进度
// @Summary 获取班级内所有视频的观看进度
// @Description 获取学生在某个班级下的所有视频学习情况
// @Tags behavior
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id path int true "班级ID"
// @Success 200 {object} http.Response{data=behavior.GetClassVideoProgressResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/behavior/class-progress/{class_id} [get]
func (bc *behaviorController) GetClassVideoProgress(c *gin.Context, req behavior.GetClassVideoProgressReq, claim ijwt.UserClaim) (http.Response, error) {
	if !domain.RoleType(claim.Role).IsValid() {
		return http.Response{}, errors.New("无权限")
	}
	items, err := bc.svc.GetClassVideoProgress(c, claim.ID, req.ClassID)
	if err != nil {
		return http.Response{}, err
	}
	out := make([]behavior.VideoProgress, 0, len(items))
	for _, it := range items {
		out = append(out, behavior.VideoProgress{
			VideoID:         it.VideoID,
			Title:           it.Title,
			Status:          it.Status,
			ProgressPercent: it.ProgressPercent,
			WatchTime:       it.WatchTime,
		})
	}
	return http.Success(behavior.GetClassVideoProgressResp{ProgressList: out}), nil
}
