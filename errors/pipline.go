package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	CreateJobErrorCode = 53000 + iota
	RunKnowledgeStageErrorCode
	RunQuizStageErrorCode
	ProcessTaskErrorCode
)

var (
	CreateJobError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, CreateJobErrorCode, "创建任务失败", err)
	}

	RunKnowledgeStageError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, RunKnowledgeStageErrorCode, "拆分知识点阶段任务失败", err)
	}

	RunQuizStageError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, RunQuizStageErrorCode, "生成习题阶段任务失败", err)
	}

	ProcessTaskError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, ProcessTaskErrorCode, "处理任务失败", err)
	}
)
