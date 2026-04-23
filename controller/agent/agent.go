package agent

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/agent"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	agentService "github.com/Knowpals/Knowpals-be-go/service/agent"
	"github.com/gin-gonic/gin"
)

type AgentController interface {
	Chat(c *gin.Context, req agent.ChatReq) (http.Response, error)
	GetChatHistory(c *gin.Context, req agent.GetChatHistoryReq) (http.Response, error)
	GetReport(c *gin.Context, req agent.GetReportReq) (http.Response, error)
	GenerateQuiz(c *gin.Context, req agent.GenerateQuizReq) (http.Response, error)
	GenerateReport(c *gin.Context, req agent.GenerateReportReq) (http.Response, error)
}

type agentController struct {
	svc agentService.AgentService
}

func NewAgentController(svc agentService.AgentService) AgentController {
	return &agentController{svc: svc}
}

// Chat 智能对话助手
// @Summary 智能对话助手
// @Tags agent
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body agent.ChatReq true "请求参数"
// @Success 200 {object} http.Response{data=agent.ChatResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/agent/chat [post]
func (ac *agentController) Chat(c *gin.Context, req agent.ChatReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := ac.svc.Chat(c, claim.ID, req)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

// GetChatHistory 获取聊天历史记录
// @Summary 获取聊天历史记录
// @Tags agent
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id query string false "视频ID（可选）"
// @Param limit query int false "条数（默认100，最大500）"
// @Success 200 {object} http.Response{data=agent.GetChatHistoryResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/agent/history [get]
func (ac *agentController) GetChatHistory(c *gin.Context, req agent.GetChatHistoryReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := ac.svc.GetChatHistory(c, claim.ID, req)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

// GetReport 获取已存储的学习报告（不触发生成）
// @Summary 获取已存储的学习报告
// @Tags agent
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id query string true "视频ID"
// @Success 200 {object} http.Response{data=agent.GenerateReportResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/agent/report [get]
func (ac *agentController) GetReport(c *gin.Context, req agent.GetReportReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := ac.svc.GenerateReport(c, claim.ID, agent.GenerateReportReq{VideoID: req.VideoID, ForceRegen: false})
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

// GenerateQuiz 生成习题
// @Summary 生成习题
// @Tags agent
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body agent.GenerateQuizReq true "请求参数"
// @Success 200 {object} http.Response{data=agent.GenerateQuizResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/agent/quiz [post]
func (ac *agentController) GenerateQuiz(c *gin.Context, req agent.GenerateQuizReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := ac.svc.GenerateQuiz(c, claim.ID, req)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

// GenerateReport 生成学情报告
// @Summary 生成学情报告
// @Tags agent
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body agent.GenerateReportReq true "请求参数"
// @Success 200 {object} http.Response{data=agent.GenerateReportResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/agent/report [post]
func (ac *agentController) GenerateReport(c *gin.Context, req agent.GenerateReportReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Student {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := ac.svc.GenerateReport(c, claim.ID, req)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

