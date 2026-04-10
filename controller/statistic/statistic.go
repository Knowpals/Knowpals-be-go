package statistic

import (
	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/statistic"
	"github.com/gin-gonic/gin"
)

type StatController interface {
	// GetStudentStat 获取学生的学习情况
	GetStudentStat(c *gin.Context, req statistic.GetStudentStatReq) (http.Response, error)
	// GetClassStat 获取班级学习情况
	GetClassStat(c *gin.Context, req statistic.GetClassStatReq) (http.Response, error)
}

type statController struct {
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
	return http.Response{}, nil
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
	return http.Response{}, nil
}
