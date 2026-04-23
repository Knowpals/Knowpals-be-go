package errors

import (
	"net/http"

	"github.com/Knowpals/Knowpals-be-go/pkg/errorx"
)

const (
	AgentChatErrorCode = 56000 + iota
	AgentGenerateQuizErrorCode
	AgentGenerateReportErrorCode
)

var (
	AgentChatError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, AgentChatErrorCode, "智能对话失败", err)
	}
	AgentGenerateQuizError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, AgentGenerateQuizErrorCode, "生成个性化题目失败", err)
	}
	AgentGenerateReportError = func(err error) error {
		return errorx.New(http.StatusInternalServerError, AgentGenerateReportErrorCode, "生成学习报告失败", err)
	}
)

