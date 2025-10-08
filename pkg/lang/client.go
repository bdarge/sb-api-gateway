package lang

import (
	"log"

	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/out/lang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClient client
type ServiceClient struct {
	Client lang.LangServiceClient
	Config config.Config
}

// InitServiceClient initialize service client
func InitServiceClient(c *config.Config) lang.LangServiceClient {
	cc, err := grpc.NewClient(c.APISvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Could not connect: %s, %v", c.APISvcURL, err)
	}

	return lang.NewLangServiceClient(cc)
}
