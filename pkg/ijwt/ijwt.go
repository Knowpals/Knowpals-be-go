package ijwt

import (
	"errors"
	"time"

	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	jwt.StandardClaims
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type JwtHandler struct {
	signingMethod jwt.SigningMethod //令牌的加密方式
	secretKey     string
	encKey        string
	timeout       int
}

func NewJwtHandler(conf *config.Config) *JwtHandler {
	return &JwtHandler{
		signingMethod: jwt.SigningMethodHS256,
		secretKey:     conf.Jwt.SecretKey,
		encKey:        conf.Jwt.EncKey,
		timeout:       conf.Jwt.Timeout,
	}
}

func (j *JwtHandler) GenerateToken(id uint, username string, password string, email string, role string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(j.timeout) * time.Second)
	claims := UserClaim{
		ID:       id,
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(j.signingMethod, claims)
	token, err := tokenClaims.SignedString([]byte(j.secretKey))
	return token, err

}

func (j *JwtHandler) ParseToken(tokenStr string) (UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		//校验签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名检验算法错误")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return UserClaim{}, err
	}

	if token == nil || !token.Valid {
		return UserClaim{}, errors.New("token无效")
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok {
		return UserClaim{}, errors.New("无法解析 token claims")
	}

	return *claims, nil
}
