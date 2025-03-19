package currency

import (
	"log"

	"github.com/bdarge/api-gateway/out/currency"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClient client
type ServiceClient struct {
	Client currency.CurrencyClient
}

// InitServiceClient init CurrencyClient
func InitServiceClient(c *config.Config) currency.CurrencyClient {
	cc, err := grpc.NewClient(c.CurrencySvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.CurrencySvcURL, err)
	}

	return currency.NewCurrencyClient(cc)
}
