package disposition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/bdarge/sb-api-gateway/pkg/utils"
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
	CreateDispositionFunc func(ctx context.Context, in *CreateDispositionRequest, opts ...grpc.CallOption) (*CreateDispositionResponse, error)
	GetDispositionFunc    func(ctx context.Context, in *GetDispositionRequest, opts ...grpc.CallOption) (*GetDispositionResponse, error)
	GetDispositionsFunc   func(ctx context.Context, in *GetDispositionsRequest, opts ...grpc.CallOption) (*GetDispositionsResponse, error)
}

func (m MockRequestServiceClient) CreateDisposition(ctx context.Context, in *CreateDispositionRequest, opts ...grpc.CallOption) (*CreateDispositionResponse, error) {
	return m.CreateDispositionFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetDisposition(ctx context.Context, in *GetDispositionRequest, opts ...grpc.CallOption) (*GetDispositionResponse, error) {
	return m.GetDispositionFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetDispositions(ctx context.Context, in *GetDispositionsRequest, opts ...grpc.CallOption) (*GetDispositionsResponse, error) {
	return m.GetDispositionsFunc(ctx, in, opts...)
}

func TestCreateDisposition(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		status       int
		order        models.Disposition
		data         *CreateDispositionResponse
	}{
		{
			name:  "should create a request",
			error: nil,
			order: models.Disposition{
				Description:  "motor",
				CreatedBy:    12341,
				CustomerId:   8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)),
				RequestType:  "order",
				Items:        []models.DispositionItem{{Description: "motor key", Qty: 1, Unit: "birr", UnitPrice: 23.4}},
			},
			status: 201,
			data:   &CreateDispositionResponse{},
		},
		{
			name:  "should return bad request #1",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "Description", Message: "This field is required"}}},
			order: models.Disposition{CreatedBy: 12341, CustomerId: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)),
				RequestType:  "order"},
			status: 400,
		},
		{
			name:  "should return bad request #2",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "CreatedBy", Message: "This field is required"}}},
			order: models.Disposition{Description: "motor", CustomerId: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "order"},
			status: 400,
		},
		{
			name:         "should return a general error when server failed to create a disposition",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			order: models.Disposition{Description: "motor", CreatedBy: 12341, CustomerId: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "order"},
			status: 500,
			data:   nil,
		},
		{
			name:  "should return bad request #3",
			error: map[string][]models.ErrorMsg{"errors": {{Field: "RequestType", Message: "Should be one of the following: 'order', or 'quote'"}}},
			order: models.Disposition{Description: "motor", CreatedBy: 12341, CustomerId: 8983,
				DeliveryDate: time.Now().Add(time.Hour * 24 * 7 * time.Duration(10)), RequestType: "cat"},
			status: 400,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.order)
		client := &MockRequestServiceClient{}
		client.CreateDispositionFunc = func(ctx context.Context, in *CreateDispositionRequest, opts ...grpc.CallOption) (*CreateDispositionResponse, error) {
			if tt.status > 201 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return tt.data, nil
			}
		}
		CreateDisposition(c, client)

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

func TestGetDisposition(t *testing.T) {
	now := time.Now()
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		data         *DispositionData
		responseData models.Disposition
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
			responseData: models.Disposition{
				Model: models.Model{
					ID:        3562,
					CreatedAt: now.UTC(),
					UpdatedAt: now.UTC(),
					DeletedAt: time.Time{}.UTC(),
				},
				CreatedBy:   30,
				CustomerId:  2,
				Description: "motor",
				RequestType: "order",
				Items: []models.DispositionItem{
					{
						Model: models.Model{
							ID:        1,
							CreatedAt: now.UTC(),
							UpdatedAt: now.UTC(),
							DeletedAt: time.Time{}.UTC(),
						},
						Description: "motor key",
						Qty:         1,
						Unit:        "birr",
						UnitPrice:   23.4,
					},
				},
				DeliveryDate: now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10)),
			},
			data: &DispositionData{
				Id:          3562,
				Description: "motor",
				RequestType: "order",
				CreatedBy:   30,
				CustomerId:  2,
				CreatedAt:   timestamppb.New(now),
				UpdatedAt:   timestamppb.New(now),
				DeletedAt:   timestamppb.New(time.Time{}),
				Items: []*DispositionItem{
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
				DeliveryDate: timestamppb.New(now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10))),
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
		client.GetDispositionFunc = func(ctx context.Context, in *GetDispositionRequest, opts ...grpc.CallOption) (*GetDispositionResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return &GetDispositionResponse{
					Data: tt.data,
				}, nil
			}
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, tt.params, url.Values{})

		GetDisposition(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			response := &models.Disposition{}
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

func TestGetDispositions(t *testing.T) {
	now := time.Now()
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		responseData models.Dispositions
		data         *GetDispositionsResponse
		status       int
		params       []gin.Param
	}{
		{
			name:   "should get dispositions",
			error:  nil,
			params: nil,
			status: 200,
			responseData: models.Dispositions{
				Limit: 10, Page: 1, Total: 1,
				Data: []models.Disposition{
					{
						Model: models.Model{
							ID:        3562,
							CreatedAt: now.UTC(),
							UpdatedAt: now.UTC(),
							DeletedAt: time.Time{}.UTC(),
						},
						Description:  "motor",
						RequestType:  "order",
						DeliveryDate: now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10)),
					},
				},
			},
			data: &GetDispositionsResponse{
				Limit: 10, Page: 1, Total: 1,
				Data: []*DispositionData{
					{
						Id:           3562,
						Description:  "motor",
						RequestType:  "order",
						CreatedAt:    timestamppb.New(now),
						UpdatedAt:    timestamppb.New(now),
						DeletedAt:    timestamppb.New(time.Time{}),
						DeliveryDate: timestamppb.New(now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10))),
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
			responseData: models.Dispositions{Limit: 10, Page: 1, Total: 1,
				Data: []models.Disposition{
					{
						Model: models.Model{
							ID:        9569,
							CreatedAt: now.UTC(),
							UpdatedAt: now.UTC(),
							DeletedAt: time.Time{}.UTC(),
						},
						Description: "motor", RequestType: "queue",
						DeliveryDate: now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10)),
					},
				},
			},
			data: &GetDispositionsResponse{
				Limit: 10, Page: 1, Total: 1,
				Data: []*DispositionData{{
					Id:          9569,
					Description: "motor", RequestType: "queue",
					CreatedAt:    timestamppb.New(now),
					UpdatedAt:    timestamppb.New(now),
					DeletedAt:    timestamppb.New(time.Time{}),
					DeliveryDate: timestamppb.New(now.UTC().Add(time.Hour * 24 * 7 * time.Duration(10))),
				}},
			},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		client := &MockRequestServiceClient{}
		client.GetDispositionsFunc = func(ctx context.Context, in *GetDispositionsRequest, opts ...grpc.CallOption) (*GetDispositionsResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			} else {
				return tt.data, nil
			}
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, tt.params, url.Values{})

		GetDispositions(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			response := &models.Dispositions{}
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
