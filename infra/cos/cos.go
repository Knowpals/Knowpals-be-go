package cos

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type COSClient struct {
	client *cos.Client
}

func NewCOSClient(cli *cos.Client) *COSClient {
	return &COSClient{
		client: cli,
	}
}

func (c *COSClient) UploadFile(
	ctx context.Context,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	}
	key := c.generateKey(fileHeader)
	_, err := c.client.Object.Put(ctx, key, file, opt)
	if err != nil {
		return "", err
	}

	//url, err := c.SignUrl(ctx, key)
	//if err != nil {
	//	return "", err
	//}
	return key, nil
}

func (c *COSClient) generateKey(fileHeader *multipart.FileHeader) string {
	//按天生成文件夹
	dayDir := time.Now().Format("20060102")
	//使用时间戳防重名
	fileName := time.Now().Format("150105") + "_" + fileHeader.Filename
	return fmt.Sprintf("%s/%s", dayDir, fileName)
}

func (c *COSClient) SignUrl(ctx context.Context, fileKey string) (string, error) {
	expire := 7 * 24 * time.Hour

	signedURL, err := c.client.Object.GetPresignedURL2(
		ctx,
		"GET",
		fileKey,
		expire,
		nil,
	)
	if err != nil {
		return "", err
	}
	return signedURL.String(), nil
}
