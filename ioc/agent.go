package ioc

import (
	"github.com/Knowpals/Knowpals-be-go/api/grpc/agentpb"
	"google.golang.org/grpc"
)

func InitAgentClient(server *grpc.ClientConn) agentpb.AgentServiceClient {
	return agentpb.NewAgentServiceClient(server)
}
