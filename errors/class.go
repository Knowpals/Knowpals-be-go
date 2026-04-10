package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	// class
	CreateClassErrorCode = 51000 + iota
	JoinClassErrorCode
	QuitClassErrorCode
	GetClassInfoErrorCode
	GetStudentClassesErrorCode
	GetTeacherClassesErrorCode
	GetClassStudentsErrorCode
)

var (
	CreateClassError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, CreateClassErrorCode, "创建班级失败", err)
	}
	JoinClassError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, JoinClassErrorCode, "加入班级失败", err)
	}
	QuitClassError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, QuitClassErrorCode, "退出班级失败", err)
	}
	GetClassInfoError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetClassInfoErrorCode, "获取班级信息失败", err)
	}
	GetStudentClassesError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetStudentClassesErrorCode, "获取学生班级列表失败", err)
	}
	GetTeacherClassesError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetTeacherClassesErrorCode, "获取教师班级列表失败", err)
	}
	GetClassStudentsError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetClassStudentsErrorCode, "获取班级学生列表失败", err)
	}
)
