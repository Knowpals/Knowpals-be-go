package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	UploadVideoErrorCode = 52000 + iota
)

var (
	UploadVideoError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, UploadVideoErrorCode, "上传视频失败", err)
	}
)
