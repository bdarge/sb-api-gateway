package disposition

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bdarge/sb-api-gateway/pkg/config"
	"github.com/bdarge/sb-api-gateway/pkg/disposition/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.DispositionServiceClient
}

func InitServiceClient(c *config.Config) pb.DispositionServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewDispositionServiceClient(cc)
}
