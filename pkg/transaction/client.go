package transaction

import (
	. "github.com/bdarge/api-gateway/out/transaction"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClient struct {
	Client TransactionServiceClient
}

func InitServiceClient(c *config.Config) TransactionServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return NewTransactionServiceClient(cc)
}
