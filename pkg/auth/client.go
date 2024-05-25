package auth

import (
	"log"

	auth "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client auth.AuthServiceClient
	Config config.Config
}

func InitServiceClient(c *config.Config) auth.AuthServiceClient {
	cc, err := grpc.Dial(c.AuthSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Could not connect: %s, %v", c.AuthSvcURL, err)
	}

	return auth.NewAuthServiceClient(cc)
}
