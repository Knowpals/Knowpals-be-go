package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Knowpals/Knowpals-be-go/api"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt *ijwt.JwtHandler
}

func NewAuthMiddleware(jwt *ijwt.JwtHandler) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

func (am *AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authStr := c.GetHeader("Authorization")
		if authStr == "" {
			c.Error(errors.New("认证头部缺失"))
			c.JSON(http.StatusUnauthorized, api.Response{
				Code:    http.StatusUnauthorized,
				Message: "认证头部缺失",
				Data:    nil,
			})
			return
		}

		segs := strings.Split(authStr, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			c.Error(errors.New("认证头部格式错误"))
			c.JSON(http.StatusUnauthorized, api.Response{
				Code:    http.StatusUnauthorized,
				Message: "认证头部格式错误",
				Data:    nil,
			})
			return
		}

		uc, err := am.jwt.ParseToken(segs[1])
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusUnauthorized, api.Response{
				Code:    http.StatusUnauthorized,
				Message: "无效或过期的身份令牌",
				Data:    nil,
			})
			return
		}

		ginx.SetClaim(c, uc)
		c.Next()

	}
}
