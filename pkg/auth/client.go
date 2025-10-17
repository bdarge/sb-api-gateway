package auth

import (
	"log"

	auth "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClient service client
type ServiceClient struct {
	Client auth.AuthServiceClient
	Config config.Config
}

// InitServiceClient initialize auth client service
func InitServiceClient(c *config.Config) auth.AuthServiceClient {
	cc, err := grpc.NewClient(c.AuthSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Could not connect: %s, %v", c.AuthSvcURL, err)
	}

	return auth.NewAuthServiceClient(cc)
}
