package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	RecordActionErrorCode = 54000 + iota
	UpdateProgressErrorCode
	GetClassVideoProgressErrorCode
)

var (
	RecordActionError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, RecordActionErrorCode, "记录视频行为失败", err)
	}

	UpdateProgressError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, UpdateProgressErrorCode, "更新观看进度失败", err)
	}

	GetClassVideoProgressError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetClassVideoProgressErrorCode, "获取班级视频进度失败", err)
	}
)
