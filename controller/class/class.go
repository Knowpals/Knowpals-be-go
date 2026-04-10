package class

import (
	errors1 "errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/class"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	class2 "github.com/Knowpals/Knowpals-be-go/service/class"
	user2 "github.com/Knowpals/Knowpals-be-go/service/user"
	"github.com/gin-gonic/gin"
)

type ClassController interface {
	// CreateClass 老师创建班级
	CreateClass(c *gin.Context, request class.CreateClassRequest, claim ijwt.UserClaim) (http.Response, error)
	// JoinClass 学生加入班级
	JoinClass(c *gin.Context, request class.JoinClassRequest, claim ijwt.UserClaim) (http.Response, error)
	// QuitClass 学生退出班级
	QuitClass(c *gin.Context, request class.QuitClassRequest, claim ijwt.UserClaim) (http.Response, error)
	// GetClassInfo 获取班级信息
	GetClassInfo(c *gin.Context, request class.GetClassInfoRequest, claim ijwt.UserClaim) (http.Response, error)

	// GetMyCreatedClasses 老师查看所有创建班级
	GetMyCreatedClasses(c *gin.Context, claim ijwt.UserClaim) (http.Response, error)
	// GetMyJoinedClasses 学生查看所有加入班级
	GetMyJoinedClasses(c *gin.Context, claim ijwt.UserClaim) (http.Response, error)
	// GetClassStudents 查看班级学生
	GetClassStudents(c *gin.Context, request class.GetClassStudentsRequest) (http.Response, error)
}

type classController struct {
	classService class2.ClassService
	userService  user2.UserService
}

func NewClassController(classService class2.ClassService, userService user2.UserService) ClassController {
	return &classController{
		classService: classService,
		userService:  userService,
	}
}

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
func (cc *classController) CreateClass(c *gin.Context, request class.CreateClassRequest, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.CreateClassError(errors1.New("无权限"))
	}

	res, err := cc.classService.CreateClass(c, claim.ID, request.ClassName)
	if err != nil {
		return http.Response{}, err
	}

	return http.Success(class.CreateClassResp{
		ClassName:  res.ClassName,
		ClassID:    res.ID,
		InviteCode: res.InviteCode,
	}), nil
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
func (cc *classController) JoinClass(c *gin.Context, request class.JoinClassRequest, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.JoinClassError(errors1.New("无权限"))
	}

	err := cc.classService.JoinClass(c, claim.ID, request.InviteCode)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
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
func (cc *classController) QuitClass(c *gin.Context, request class.QuitClassRequest, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.QuitClassError(errors1.New("无权限"))
	}
	err := cc.classService.QuitClass(c, claim.ID, request.ClassID)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
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
func (cc *classController) GetClassInfo(c *gin.Context, request class.GetClassInfoRequest, claim ijwt.UserClaim) (http.Response, error) {
	if !domain.RoleType(claim.Role).IsValid() {
		return http.Response{}, errors.GetClassInfoError(errors1.New("无权限"))
	}

	classDomain, err := cc.classService.GetClassByID(c, request.ClassID)
	if err != nil {
		return http.Response{}, err
	}

	teacher, err := cc.userService.GetUserByID(c, classDomain.TeacherID)
	if err != nil {
		return http.Response{}, err
	}

	return http.Success(class.GetClassInfoResp{
		ClassInfo: class.ClassInfo{
			TeacherID:   classDomain.TeacherID,
			TeacherName: teacher.Username,
			ClassID:     classDomain.ID,
			ClassName:   classDomain.ClassName,
			InviteCode:  classDomain.InviteCode,
		},
	}), nil
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
func (cc *classController) GetMyCreatedClasses(c *gin.Context, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.GetTeacherClassesError(errors1.New("无权限"))
	}

	classes, err := cc.classService.GetTeacherClasses(c, claim.ID)
	if err != nil {
		return http.Response{}, err
	}

	teacherName := claim.Username
	classInfos := make([]class.ClassInfo, 0, len(classes))
	for _, cls := range classes {
		classInfos = append(classInfos, class.ClassInfo{
			TeacherID:   cls.TeacherID,
			TeacherName: teacherName,
			ClassID:     cls.ID,
			ClassName:   cls.ClassName,
			InviteCode:  cls.InviteCode,
		})
	}
	return http.Success(class.GetMyCreatedClassesResp{ClassList: classInfos}), nil
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
func (cc *classController) GetMyJoinedClasses(c *gin.Context, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.GetStudentClassesError(errors1.New("无权限"))
	}

	classes, err := cc.classService.GetStudentClasses(c, claim.ID)
	if err != nil {
		return http.Response{}, err
	}

	teacherNameCache := make(map[uint]string, 8)
	getTeacherName := func(teacherID uint) (string, error) {
		if name, ok := teacherNameCache[teacherID]; ok {
			return name, nil
		}
		u, err := cc.userService.GetUserByID(c, teacherID)
		if err != nil {
			return "", err
		}
		teacherNameCache[teacherID] = u.Username
		return u.Username, nil
	}

	classInfos := make([]class.ClassInfo, 0, len(classes))
	for _, cls := range classes {
		name, err := getTeacherName(cls.TeacherID)
		if err != nil {
			return http.Response{}, err
		}
		classInfos = append(classInfos, class.ClassInfo{
			TeacherID:   cls.TeacherID,
			TeacherName: name,
			ClassID:     cls.ID,
			ClassName:   cls.ClassName,
			InviteCode:  cls.InviteCode,
		})
	}

	return http.Success(class.GetMyJoinedClassesResp{ClassList: classInfos}), nil
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

	students, err := cc.classService.GetClassStudents(c, request.ClassID)
	if err != nil {
		return http.Response{}, err
	}

	resp := class.GetClassStudentsResp{
		Students: make([]class.Student, 0, len(students)),
	}
	for _, s := range students {
		if domain.RoleType(s.Role) != domain.Role_Student {
			continue
		}
		resp.Students = append(resp.Students, class.Student{
			ID:       s.ID,
			Username: s.Username,
			Email:    s.Email,
		})
	}

	return http.Success(resp), nil
}
