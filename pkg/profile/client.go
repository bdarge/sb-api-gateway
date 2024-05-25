package profile

import (
	"log"

	. "github.com/bdarge/api-gateway/out/profile"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client ProfileServiceClient
}

func InitServiceClient(c *config.Config) ProfileServiceClient {
	cc, err := grpc.Dial(c.APISvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.APISvcURL, err)
	}

	return NewProfileServiceClient(cc)
}
