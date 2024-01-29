package transaction

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/bdarge/api-gateway/out/model"
	. "github.com/bdarge/api-gateway/out/transaction"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type MockRequestServiceClient struct {
	CreateTransactionFunc  func(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*CreateTransactionResponse, error)
	GetTransactionFunc     func(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	GetTransactionsFunc    func(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error)
	UpdateTransactionsFunc func(ctx context.Context, in *UpdateTransactionRequest, opts ...grpc.CallOption) (*UpdateTransactionResponse, error)
	DeleteTransactionsFunc func(ctx context.Context, in *DeleteTransactionRequest, opts ...grpc.CallOption) (*DeleteTransactionResponse, error)
}

func (m MockRequestServiceClient) CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*CreateTransactionResponse, error) {
	return m.CreateTransactionFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
	return m.GetTransactionFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error) {
	return m.GetTransactionsFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) UpdateTransaction(ctx context.Context, in *UpdateTransactionRequest, opts ...grpc.CallOption) (*UpdateTransactionResponse, error) {
	return m.UpdateTransactionsFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) DeleteTransaction(ctx context.Context, in *DeleteTransactionRequest, opts ...grpc.CallOption) (*DeleteTransactionResponse, error) {
	return m.DeleteTransactionsFunc(ctx, in, opts...)
}

func TestCreateTransaction(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		status       int
		order        models.Transaction
		data         *CreateTransactionResponse
	}{
		{
			name:  "should create a request",
			error: nil,
			order: models.Transaction{
				Description:  "motor",
				CreatedBy:    12341,
				CustomerID:   8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)),
				RequestType:  "order",
				Items:        []models.TransactionItem{{Description: "motor key", Qty: "1", Unit: "birr", UnitPrice: "23.4"}},
			},
			status: 201,
			data:   &CreateTransactionResponse{},
		},
		{
			name:  "should return bad request #1",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "Description", Message: "This field is required"}}},
			order: models.Transaction{CreatedBy: 12341, CustomerID: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)),
				RequestType:  "order"},
			status: 400,
		},
		{
			name:  "should return bad request #2",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "CreatedBy", Message: "This field is required"}}},
			order: models.Transaction{Description: "motor", CustomerID: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "order"},
			status: 400,
		},
		{
			name:         "should return a general error when server failed to create a disposition",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			order: models.Transaction{Description: "motor", CreatedBy: 12341, CustomerID: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "order"},
			status: 500,
			data:   nil,
		},
		{
			name:  "should return bad request #3",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "RequestType", Message: "Should be one of the following: 'order', or 'quote'"}}},
			order: models.Transaction{Description: "motor", CreatedBy: 12341, CustomerID: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "cat"},
			status: 400,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.order)
		client := &MockRequestServiceClient{}
		client.CreateTransactionFunc = func(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*CreateTransactionResponse, error) {
			if tt.status > 201 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return tt.data, nil
			}
		}
		CreateTransaction(c, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code > 201 {
			b, _ := io.ReadAll(w.Body)
			if tt.error != nil {
				var targetedError map[string][]models.ErrorMsg
				if err := json.Unmarshal(b, &targetedError); err != nil {
					t.Error(tt.name, "invalid error type", string(b))
				}
				if !reflect.DeepEqual(targetedError, tt.error) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			} else if tt.generalError != nil {
				var generalError map[string]string
				if err := json.Unmarshal(b, &generalError); err != nil {
					if err := json.Unmarshal(b, &generalError); err != nil {
						t.Error(tt.name, "invalid error type", string(b))
					}
				}
				if !reflect.DeepEqual(generalError, tt.generalError) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			}
		}
	}
}

func TestGetTransaction(t *testing.T) {
	now := time.Now()
	nowInUtc := &now
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		data         *TransactionData
		responseData models.Transaction
		status       int
		params       []gin.Param
	}{
		{
			name:  "should get a disposition:",
			error: nil,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "37623",
				},
			},
			status: 200,
			responseData: models.Transaction{
				Model: models.Model{
					ID:        3562,
					CreatedAt: nowInUtc,
					UpdatedAt: nowInUtc,
					DeletedAt: nil,
				},
				CreatedBy:   30,
				CustomerID:  2,
				Description: "motor",
				RequestType: "order",
				Items: []models.TransactionItem{
					{
						Model: models.Model{
							ID:        1,
							CreatedAt: nowInUtc,
							UpdatedAt: nowInUtc,
							DeletedAt: nil,
						},
						Description: "motor key",
						Qty:         "1",
						Unit:        "birr",
						UnitPrice:   "23.4",
					},
				},
				DeliveryDate: nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10)),
			},
			data: &TransactionData{
				Id:          3562,
				Description: "motor",
				RequestType: "order",
				CreatedBy:   30,
				CustomerId:  2,
				CreatedAt:   timestamppb.New(now),
				UpdatedAt:   timestamppb.New(now),
				DeletedAt:   timestamppb.New(time.Time{}),
				Items: []*TransactionItem{
					{
						Id:          1,
						Description: "motor key",
						Qty:         1,
						Unit:        "birr",
						UnitPrice:   23.4,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
						DeletedAt:   timestamppb.New(time.Time{}),
					},
				},
				DeliveryDate: timestamppb.New(nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10))),
			},
		},
		{
			name:         "should return an error:",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			params: []gin.Param{
				{
					Key:   "id",
					Value: "647364",
				},
			},
			status: 500,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		client := &MockRequestServiceClient{}
		client.GetTransactionFunc = func(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return &GetTransactionResponse{
					Data: tt.data,
				}, nil
			}
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, tt.params, url.Values{})

		GetTransaction(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			response := &models.Transaction{}
			err := json.Unmarshal(b, response)
			if err != nil {
				t.Error(tt.name, "test error", err)
				continue
			}

			if !reflect.DeepEqual(response, &tt.responseData) {
				t.Error(tt.name, "data doesn't match,", "expected:", tt.responseData, "actual:", string(b))
			}
		}

		if w.Code > 200 {
			b, _ := io.ReadAll(w.Body)
			if tt.error != nil {
				var targetedError map[string][]models.ErrorMsg
				if err := json.Unmarshal(b, &targetedError); err != nil {
					t.Error(tt.name, "invalid error type", string(b))
				}
				if !reflect.DeepEqual(targetedError, tt.error) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			} else if tt.generalError != nil {
				var generalError map[string]string
				if err := json.Unmarshal(b, &generalError); err != nil {
					if err := json.Unmarshal(b, &generalError); err != nil {
						t.Error(tt.name, "invalid error type", string(b))
					}
				}
				if !reflect.DeepEqual(generalError, tt.generalError) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			}
		}
	}
}

func TestGetTransactions(t *testing.T) {
	now := time.Now()
	nowInUtc := &now
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		responseData models.Transactions
		data         *GetTransactionsResponse
		status       int
		params       []gin.Param
	}{
		{
			name:   "should get dispositions",
			error:  nil,
			params: nil,
			status: 200,
			responseData: models.Transactions{
				Limit: 10, Page: 1, Total: 1,
				Data: []models.Transaction{
					{
						Model: models.Model{
							ID:        3562,
							CreatedAt: nowInUtc,
							UpdatedAt: nowInUtc,
							DeletedAt: nil,
						},
						Description:  "motor",
						RequestType:  "order",
						DeliveryDate: nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10)),
					},
				},
			},
			data: &GetTransactionsResponse{
				Limit: 10, Page: 1, Total: 1,
				Data: []*TransactionData{
					{
						Id:           3562,
						Description:  "motor",
						RequestType:  "order",
						CreatedAt:    timestamppb.New(now),
						UpdatedAt:    timestamppb.New(now),
						DeletedAt:    timestamppb.New(time.Time{}),
						DeliveryDate: timestamppb.New(nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10))),
					},
				}},
		},
		{
			name:  "should get dispositions by requestType if requestType is sent",
			error: nil,
			params: []gin.Param{
				{
					Key:   "requestType",
					Value: "order",
				},
			},
			status: 200,
			responseData: models.Transactions{Limit: 10, Page: 1, Total: 1,
				Data: []models.Transaction{
					{
						Model: models.Model{
							ID:        9569,
							CreatedAt: nowInUtc,
							UpdatedAt: nowInUtc,
							DeletedAt: nil,
						},
						Description: "motor", RequestType: "queue",
						DeliveryDate: nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10)),
					},
				},
			},
			data: &GetTransactionsResponse{
				Limit: 10, Page: 1, Total: 1,
				Data: []*TransactionData{{
					Id:          9569,
					Description: "motor", RequestType: "queue",
					CreatedAt:    timestamppb.New(now),
					UpdatedAt:    timestamppb.New(now),
					DeletedAt:    timestamppb.New(time.Time{}),
					DeliveryDate: timestamppb.New(nowInUtc.Add(time.Hour * 24 * 7 * time.Duration(10))),
				}},
			},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		client := &MockRequestServiceClient{}
		client.GetTransactionsFunc = func(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return tt.data, nil
			}
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, tt.params, url.Values{})

		GetTransactions(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			response := &models.Transactions{}
			err := json.Unmarshal(b, response)
			if err != nil {
				t.Error(tt.name, "test error", err)
				continue
			}
			if !reflect.DeepEqual(*response, tt.responseData) {
				t.Error(tt.name, "data doesn't match, actual:", *response, "Expected:", tt.responseData)
			}
		}

		if w.Code > 200 {
			b, _ := io.ReadAll(w.Body)
			if tt.error != nil {
				var targetedError map[string][]models.ErrorMsg
				if err := json.Unmarshal(b, &targetedError); err != nil {
					t.Error(tt.name, "invalid error type", string(b))
				}
				if !reflect.DeepEqual(targetedError, tt.error) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			} else if tt.generalError != nil {
				var generalError map[string]string
				if err := json.Unmarshal(b, &generalError); err != nil {
					if err := json.Unmarshal(b, &generalError); err != nil {
						t.Error(tt.name, "invalid error type", string(b))
					}
				}
				if !reflect.DeepEqual(generalError, tt.generalError) {
					t.Error(tt.name, "error doesn't match,", string(b))
				}
			}
		}
	}
}
