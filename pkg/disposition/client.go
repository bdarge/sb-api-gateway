package disposition

import (
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/bdarge/sb-api-gateway/pkg/config"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client DispositionServiceClient
}

func InitServiceClient(c *config.Config) DispositionServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return NewDispositionServiceClient(cc)
}
