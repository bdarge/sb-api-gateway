package customer

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/bdarge/api-gateway/out/customer"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type MockRequestServiceClient struct {
	CreatCustomerFunc func(ctx context.Context, in *CreateCustomerRequest, opts ...grpc.CallOption) (*CreateCustomerResponse, error)
	GetCustomerFunc   func(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerResponse, error)
}

func (m MockRequestServiceClient) CreateCustomer(ctx context.Context, in *CreateCustomerRequest, opts ...grpc.CallOption) (*CreateCustomerResponse, error) {
	return m.CreatCustomerFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerResponse, error) {
	return m.GetCustomerFunc(ctx, in, opts...)
}

func TestCreateCustomer(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		status       int
		order        models.Customer
		data         *CreateCustomerResponse
	}{
		{
			name:  "should create a customer",
			error: nil,
			order: models.Customer{
				Name:  "John Ali",
				Email: "fake@gmail.com",
			},
			status: 201,
			data:   &CreateCustomerResponse{},
		},
		{
			name:  "should return bad request #1",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "Email", Message: "This field is required"}}},
			order: models.Customer{
				Name: "John Ali",
			},
			status: 400,
		},
		{
			name:  "should return bad request #1",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "Name", Message: "This field is required"}}},
			order: models.Customer{
				Email: "John Ali",
			},
			status: 400,
		},
		{
			name:         "should return a general error when server failed to create a Customer",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			order:        models.Customer{Name: "Mike", Email: "fake@gmail.com"},
			status:       500,
			data:         nil,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.order)
		client := &MockRequestServiceClient{}
		client.CreatCustomerFunc = func(ctx context.Context, in *CreateCustomerRequest, opts ...grpc.CallOption) (*CreateCustomerResponse, error) {
			if tt.status > 201 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return tt.data, nil
			}
		}

		// act
		CreateCustomer(c, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, "with status code: ", w.Code, string(b))
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

func TestGetCustomer(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		data         *CustomerData
		status       int
		params       []gin.Param
	}{
		{
			name:  "should get a customer",
			error: nil,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "37623",
				},
			},
			status: 200,
			data: &CustomerData{
				Name:  "Mike Teddy",
				Email: "fake@gmail.com",
			},
		},
		{
			name:         "should return an error if there is an internal issue",
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
		client.GetCustomerFunc = func(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return &GetCustomerResponse{
					Data: tt.data,
				}, nil
			}
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, tt.params, url.Values{})

		// act
		GetCustomer(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			d := &CustomerData{}
			err := json.Unmarshal(b, d)
			if err != nil {
				t.Error(tt.name, "test error:", err)
				continue
			}
			if !reflect.DeepEqual(d, tt.data) {
				t.Error(tt.name, "data doesn't match actual:", string(b), "expected:", tt.data)
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
