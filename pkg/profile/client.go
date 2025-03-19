package profile

import (
	"log"

	. "github.com/bdarge/api-gateway/out/profile"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClient client struct
type ServiceClient struct {
	Client ProfileServiceClient
}

// InitServiceClient initialize profile service client
func InitServiceClient(c *config.Config) ProfileServiceClient {
	cc, err := grpc.NewClient(c.APISvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.APISvcURL, err)
	}

	return NewProfileServiceClient(cc)
}
