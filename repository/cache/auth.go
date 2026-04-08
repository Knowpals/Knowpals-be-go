package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const VerifyCodePrefix = "knowpals:auth:code:"

type AuthCache interface {
	SetCode(ctx context.Context, code string, email string) error
	GetCode(ctx context.Context, email string) (string, error)
}

type authCache struct {
	cmd *redis.Client
}

func NewAuthCache(cmd *redis.Client) AuthCache {
	return &authCache{cmd: cmd}
}

func (ac *authCache) SetCode(ctx context.Context, code string, email string) error {
	//如果是同一个人验证码直接覆盖，保证只有一个有效验证码
	key := ac.getKey(email)
	return ac.cmd.Set(ctx, key, code, 5*time.Minute).Err()
}

func (ac *authCache) GetCode(ctx context.Context, email string) (string, error) {
	key := ac.getKey(email)
	code, err := ac.cmd.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

func (ac *authCache) getKey(email string) string {
	return VerifyCodePrefix + email
}
