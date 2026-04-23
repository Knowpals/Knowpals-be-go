package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	UploadVideoErrorCode = 52000 + iota
	GetVideoDetailErrorCode
	VideoReviewErrorCode
	VideoPublishErrorCode
	VideoNotPublishedErrorCode
)

var (
	UploadVideoError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, UploadVideoErrorCode, "上传视频失败", err)
	}

	GetVideoDetailError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, GetVideoDetailErrorCode, "获取视频详情失败", err)
	}

	VideoReviewError = func(err error) error {
		return errorx.New(http.StatusBadRequest, VideoReviewErrorCode, "视频进入审核失败", err)
	}

	VideoPublishError = func(err error) error {
		return errorx.New(http.StatusBadRequest, VideoPublishErrorCode, "视频发布失败", err)
	}

	VideoNotPublishedError = func(err error) error {
		return errorx.New(http.StatusForbidden, VideoNotPublishedErrorCode, "视频未发布", err)
	}
)
