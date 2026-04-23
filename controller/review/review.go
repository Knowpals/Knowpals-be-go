package review

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	httpReview "github.com/Knowpals/Knowpals-be-go/api/http/review"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	videoSvc "github.com/Knowpals/Knowpals-be-go/service/video"
	"github.com/gin-gonic/gin"
)

type ReviewController interface {
	Start(c *gin.Context, req httpReview.StartReviewReq) (http.Response, error)
	Publish(c *gin.Context, req httpReview.PublishReq) (http.Response, error)
}

type reviewController struct {
	vs videoSvc.VideoService
}

func NewReviewController(vs videoSvc.VideoService) ReviewController {
	return &reviewController{vs: vs}
}

// Start 视频进入审核
// @Summary 视频进入审核
// @Tags review
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id path int true "视频ID"
// @Success 200 {object} http.Response "成功"
// @Router /api/v1/video/{video_id}/review/start [post]
func (rc *reviewController) Start(c *gin.Context, req httpReview.StartReviewReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	if err := rc.vs.StartReview(c, req.VideoID); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// Publish 发布视频（审核完成）
// @Summary 发布视频
// @Tags review
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id path int true "视频ID"
// @Success 200 {object} http.Response "成功"
// @Router /api/v1/video/{video_id}/review/publish [post]
func (rc *reviewController) Publish(c *gin.Context, req httpReview.PublishReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	if err := rc.vs.Publish(c, req.VideoID); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

