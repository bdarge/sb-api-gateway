package auth

import (
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/bdarge/sb-api-gateway/pkg/config"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client AuthServiceClient
}

func InitServiceClient(c *config.Config) AuthServiceClient {
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Could not connect: %s, %v", c.AuthSvcUrl, err)
	}

	return NewAuthServiceClient(cc)
}
