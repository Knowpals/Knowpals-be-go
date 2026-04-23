package ioc

import (
	"context"
	"time"

	"github.com/Knowpals/Knowpals-be-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitGRPCConn 直连 Python gRPC（默认 localhost:50051）
// 注意：本项目不使用 etcd 做服务发现。
func InitGRPCConn(conf *config.Config) *grpc.ClientConn {
	addr := conf.Grpc.Addr
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}
	return conn
}
