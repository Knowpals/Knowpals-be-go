package question

import (
	"errors"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	httpQuestion "github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/gin-gonic/gin"
)

// ReviewListQuestions 老师查看待审核题目（pipeline 公共题）
func (qc *questionController) ReviewListQuestions(c *gin.Context, req httpQuestion.ReviewListReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	resp, err := qc.svc.ReviewListQuestions(c, claim.ID, domain.RoleType(claim.Role), req.VideoID)
	if err != nil {
		return http.Response{}, err
	}
	return http.Success(resp), nil
}

// ReviewAddQuestion 老师按时间点添加题目
func (qc *questionController) ReviewAddQuestion(c *gin.Context, req httpQuestion.ReviewAddReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	if err := qc.svc.ReviewAddQuestion(c, claim.ID, domain.RoleType(claim.Role), req); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// ReviewUpdateQuestion 老师修改题目（不改 segment）
func (qc *questionController) ReviewUpdateQuestion(c *gin.Context, req httpQuestion.ReviewUpdateReq) (http.Response, error) {
	// 该接口需要同时支持 uri(question_id) + json body，这里补一次 BindJSON
	_ = c.ShouldBindJSON(&req)

	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	if err := qc.svc.ReviewUpdateQuestion(c, claim.ID, domain.RoleType(claim.Role), req); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// ReviewDeleteQuestion 老师删除题目
func (qc *questionController) ReviewDeleteQuestion(c *gin.Context, req httpQuestion.ReviewDeleteReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, err
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.New("无权限")
	}
	if err := qc.svc.ReviewDeleteQuestion(c, claim.ID, domain.RoleType(claim.Role), req.QuestionID); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}
