package class

import (
	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/class"
	"github.com/gin-gonic/gin"
)

type ClassController interface {
	// CreateClass 老师创建班级
	CreateClass(c *gin.Context, request class.CreateClassRequest) (http.Response, error)
	// JoinClass 学生加入班级
	JoinClass(c *gin.Context, request class.JoinClassRequest) (http.Response, error)
	// QuitClass 学生退出班级
	QuitClass(c *gin.Context, request class.QuitClassRequest) (http.Response, error)
	// GetClassInfo 获取班级信息
	GetClassInfo(c *gin.Context, request class.GetClassInfoRequest) (http.Response, error)

	// GetMyCreatedClasses 老师查看所有创建班级
	GetMyCreatedClasses(c *gin.Context) (http.Response, error)
	// GetMyJoinedClasses 学生查看所有加入班级
	GetMyJoinedClasses(c *gin.Context) (http.Response, error)
	// GetClassStudents 查看班级学生
	GetClassStudents(c *gin.Context, request class.GetClassStudentsRequest) (http.Response, error)
}

type classController struct{}

// CreateClass 创建班级
// @Summary 教师创建班级
// @Description 教师创建班级，返回班级信息和邀请码
// @Tags class
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body class.CreateClassRequest true "创建班级参数"
// @Success 200 {object} http.Response{data=class.CreateClassResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/class/create [post]
func (cc *classController) CreateClass(c *gin.Context, request class.CreateClassRequest) (http.Response, error) {
	return http.Response{}, nil
}

// JoinClass 加入班级
// @Summary 学生加入班级
// @Description 学生使用邀请码加入班级
// @Tags class
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body class.JoinClassRequest true "邀请码"
// @Success 200 {object} http.Response "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/class/join [post]
func (cc *classController) JoinClass(c *gin.Context, request class.JoinClassRequest) (http.Response, error) {
	return http.Response{}, nil
}

// QuitClass 退出班级
// @Summary 学生退出班级
// @Description 学生退出指定班级，单参数使用 path
// @Tags class
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id path int true "班级ID"
// @Success 200 {object} http.Response "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/class/quit/{class_id} [post]
func (cc *classController) QuitClass(c *gin.Context, request class.QuitClassRequest) (http.Response, error) {
	return http.Response{}, nil
}

// GetClassInfo 获取班级信息
// @Summary 获取班级详情
// @Description 获取班级基本信息，单参数使用 path
// @Tags class
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id path int true "班级ID"
// @Success 200 {object} http.Response{data=class.GetClassInfoResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/class/info/{class_id} [get]
func (cc *classController) GetClassInfo(c *gin.Context, request class.GetClassInfoRequest) (http.Response, error) {
	return http.Response{}, nil
}

// GetMyCreatedClasses 我创建的班级
// @Summary 获取教师创建的班级列表
// @Description 获取当前老师创建的所有班级
// @Tags class
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} http.Response{data=class.GetMyCreatedClassesResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/class/my-created [get]
func (cc *classController) GetMyCreatedClasses(c *gin.Context) (http.Response, error) {
	return http.Response{}, nil
}

// GetMyJoinedClasses 我加入的班级
// @Summary 获取学生加入的班级列表
// @Description 获取当前学生加入的所有班级
// @Tags class
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} http.Response{data=class.GetMyJoinedClassesResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/class/my-joined [get]
func (cc *classController) GetMyJoinedClasses(c *gin.Context) (http.Response, error) {
	return http.Response{}, nil
}

// GetClassStudents 获取班级学生
// @Summary 获取班级内学生列表
// @Description 获取指定班级的所有学生，单参数使用 path
// @Tags class
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id path int true "班级ID"
// @Success 200 {object} http.Response{data=class.GetClassStudentsResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Failure 400 {object} http.Response "参数错误"
// @Router /api/v1/class/students/{class_id} [get]
func (cc *classController) GetClassStudents(c *gin.Context, request class.GetClassStudentsRequest) (http.Response, error) {
	return http.Response{}, nil
}
