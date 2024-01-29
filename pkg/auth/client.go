package auth

import (
	auth "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClient struct {
	Client auth.AuthServiceClient
	Config config.Config
}

func InitServiceClient(c *config.Config) auth.AuthServiceClient {
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Could not connect: %s, %v", c.AuthSvcUrl, err)
	}

	return auth.NewAuthServiceClient(cc)
}
