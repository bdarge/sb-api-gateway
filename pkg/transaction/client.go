package transaction

import (
	"github.com/bdarge/api-gateway/out/transaction"
	"github.com/bdarge/api-gateway/out/transactionItem"
	"github.com/bdarge/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// ServiceClient service client struct
type ServiceClient struct {
	Client transaction.TransactionServiceClient
}

// InitServiceClient initalize ServiceClient
func InitServiceClient(c *config.Config) transaction.TransactionServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return transaction.NewTransactionServiceClient(cc)
}

// TranItemServiceClient service client struct for TransactionItem
type TranItemServiceClient struct {
	Client transactionItem.TransactionItemServiceClient
}

// InitTranItemServiceClient initalize TransactionItemServiceClient
func InitTranItemServiceClient(c *config.Config) transactionItem.TransactionItemServiceClient {
	cc, err := grpc.Dial(c.ApiSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("couldn't connect to %s: %s", c.ApiSvcUrl, err)
	}

	return transactionItem.NewTransactionItemServiceClient(cc)
}

