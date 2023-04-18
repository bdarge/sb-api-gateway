package profile

import (
	. "github.com/bdarge/api-gateway/out/profile"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClient struct {
	Client ProfileServiceClient
}

func InitServiceClient(c *config.Config) ProfileServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return NewProfileServiceClient(cc)
}
