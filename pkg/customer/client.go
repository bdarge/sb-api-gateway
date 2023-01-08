package customer

import (
	"github.com/bdarge/sb-api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClient struct {
	Client CustomerServiceClient
}

func InitServiceClient(c *config.Config) CustomerServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return NewCustomerServiceClient(cc)
}
