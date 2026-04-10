package behavior

import (
	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/behavior"
	"github.com/gin-gonic/gin"
)

type BehaviorController interface {
	// RecordAction 记录学生观看行为
	RecordAction(c *gin.Context, req behavior.RecordActionReq) (http.Response, error)
	// UpdateVideoProgress 更新学生观看进度
	UpdateVideoProgress(c *gin.Context, req behavior.UpdateProgressReq) (http.Response, error)
	// GetClassVideoProgress 获取学生班级内所有视频观看进度
	GetClassVideoProgress(c *gin.Context, req behavior.GetClassVideoProgressReq) (http.Response, error)
}

type behaviorController struct{}

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
func (bc *behaviorController) RecordAction(c *gin.Context, req behavior.RecordActionReq) (http.Response, error) {
	return http.Response{}, nil
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
func (bc *behaviorController) UpdateVideoProgress(c *gin.Context, req behavior.UpdateProgressReq) (http.Response, error) {
	return http.Response{}, nil
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
func (bc *behaviorController) GetClassVideoProgress(c *gin.Context, req behavior.GetClassVideoProgressReq) (http.Response, error) {
	return http.Response{}, nil
}
