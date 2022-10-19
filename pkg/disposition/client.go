package request

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bdarge/sb-api-gateway/pkg/config"
	"github.com/bdarge/sb-api-gateway/pkg/request/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.RequestServiceClient
}

func InitServiceClient(c *config.Config) pb.RequestServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewRequestServiceClient(cc)
}
