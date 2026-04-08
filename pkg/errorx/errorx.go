package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type CustomError struct {
	HttpCode int    `json:"http_code"` //http错误
	Code     int    `json:"code"`      //具体业务错误码
	Msg      string `json:"message"`   //暴露给前端的错误信息
	Err      error  `json:"_"`         //具体错误原因（不暴露）
	File     string `json:"-"`         //出错文件名（不暴露）
	Line     int    `json:"-"`         //出错行号（不暴露）
	Function string `json:"-"`         //出错函数名（不暴露）

}

// 这个错误是给内部看的，传给前端的只有暴露的部分
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s (at %s:%d in %s): %v", e.Code, e.Msg, e.File, e.Line, e.Function, e.Err)
	}
	return fmt.Sprintf("[%d] %s (at %s:%d in %s)", e.Code, e.Msg, e.File, e.Line, e.Function)

}

func New(httpCode int, code int, message string, err error) error {
	//获取调用栈信息
	file, line, function := getCallerInfo(3)
	return &CustomError{
		HttpCode: httpCode,
		Code:     code,
		Msg:      message,
		Err:      err,
		File:     file,
		Line:     line,
		Function: function,
	}
}

// getCallerInfo 获取调用信息
func getCallerInfo(skip int) (string, int, string) {
	// skip: 调用栈层级，1 表示当前函数，2 表示上层调用函数 3表示上上级函数
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}
	function := runtime.FuncForPC(pc).Name()
	return file, line, function

}

func ToCustomError(err error) *CustomError {
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return customErr
	}
	file, line, function := getCallerInfo(4)
	return &CustomError{
		HttpCode: http.StatusInternalServerError,
		Code:     50001, // 通用内部错误码
		Msg:      "服务器内部错误",
		Err:      err,
		File:     file,
		Line:     line,
		Function: function,
	}

}
