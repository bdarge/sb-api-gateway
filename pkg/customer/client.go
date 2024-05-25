package customer

import (
	"log"

	"github.com/bdarge/api-gateway/out/customer"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClient client
type ServiceClient struct {
	Client customer.CustomerServiceClient
}

// InitServiceClient init CustomerServiceClient
func InitServiceClient(c *config.Config) customer.CustomerServiceClient {
	cc, err := grpc.Dial(c.APISvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.APISvcURL, err)
	}

	return customer.NewCustomerServiceClient(cc)
}
