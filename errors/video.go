package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	UploadVideoErrorCode = 52000 + iota
	GetVideoDetailErrorCode
)

var (
	UploadVideoError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, UploadVideoErrorCode, "上传视频失败", err)
	}

	GetVideoDetailError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetVideoDetailErrorCode, "获取视频详情失败", err)
	}
)
