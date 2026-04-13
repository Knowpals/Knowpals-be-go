package ioc

import (
	"net/http"
	"net/url"

	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func InitCOS(conf *config.Config) *cos.Client {
	u, _ := url.Parse(conf.COS.URL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.COS.SecretID,
			SecretKey: conf.COS.SecretKey,
		},
	})
	return client
}
