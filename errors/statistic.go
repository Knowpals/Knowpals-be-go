package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	GetStudentStatErrorCode = 55000 + iota
	GetClassStatErrorCode
	GetStudentOverviewErrorCode
)

var (
	GetStudentStatError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetStudentStatErrorCode, "获取学生统计情况", err)
	}

	GetClassStatError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetClassStatErrorCode, "获取班级统计情况", err)
	}

	GetStudentOverviewError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetStudentOverviewErrorCode, "获取学生总览统计情况", err)
	}
)
