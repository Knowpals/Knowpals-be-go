package ioc

import (
	"github.com/Knowpals/Knowpals-be-go/api/grpc/memorypb"
	"google.golang.org/grpc"
)

func InitMemoryClient(server *grpc.ClientConn) memorypb.MemoryServiceClient {
	return memorypb.NewMemoryServiceClient(server)
}
