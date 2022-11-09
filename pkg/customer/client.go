package customer

import (
	"github.com/bdarge/sb-api-gateway/pkg/config"
	"github.com/bdarge/sb-api-gateway/pkg/customer/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClient struct {
	Client pb.CustomerServiceClient
}

func InitServiceClient(c *config.Config) pb.CustomerServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return pb.NewCustomerServiceClient(cc)
}
