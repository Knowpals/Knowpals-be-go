package email

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	cmd    *redis.Client
	server string
	secret string
	addr   string
}

func NewEmailClient(cmd *redis.Client, conf *config.Config) *EmailClient {
	return &EmailClient{
		cmd:    cmd,
		secret: conf.Smtp.Secret,
		server: conf.Smtp.Server,
		addr:   conf.Smtp.Addr,
	}
}

func (e *EmailClient) SendEmail(target string, text string) error {
	m := gomail.NewMessage()

	// 发件人
	m.SetHeader("From", e.addr)
	// 收件人
	m.SetHeader("To", target)
	// 标题
	m.SetHeader("Subject", "Verification Code")
	// 正文
	m.SetBody("text/plain", text)

	// QQ 邮箱 SSL 端口 465 专用写法
	d := gomail.NewDialer(
		"smtp.qq.com",
		465,
		e.addr,
		e.secret,
	)

	// 发送
	return d.DialAndSend(m)
}
