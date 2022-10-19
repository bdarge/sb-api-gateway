package request

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/bdarge/sb-api-gateway/pkg/request/pb"
	"github.com/bdarge/sb-api-gateway/pkg/utils"
	"google.golang.org/grpc"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockRequestServiceClient struct {
	CreateRequestFunc func(ctx context.Context, in *pb.CreateRequestRequest, opts ...grpc.CallOption) (*pb.CreateRequestResponse, error)
	GetRequestFunc    func(ctx context.Context, in *pb.GetRequestRequest, opts ...grpc.CallOption) (*pb.GetRequestResponse, error)
}

func (m MockRequestServiceClient) CreateRequest(ctx context.Context, in *pb.CreateRequestRequest, opts ...grpc.CallOption) (*pb.CreateRequestResponse, error) {
	return m.CreateRequestFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetRequest(ctx context.Context, in *pb.GetRequestRequest, opts ...grpc.CallOption) (*pb.GetRequestResponse, error) {
	return m.GetRequestFunc(ctx, in, opts...)
}

func TestCreateRequest(t *testing.T) {
	// test cases
	tests := []struct {
		name   string
		error  map[string]string
		status int
		order  models.Request
		server func(ctx context.Context, in *pb.CreateRequestRequest, opts ...grpc.CallOption) (*pb.CreateRequestResponse, error)
	}{
		{
			name:   "should create an request",
			error:  nil,
			order:  models.Request{Description: "motor", CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022"},
			status: 201,
			server: func(ctx context.Context, in *pb.CreateRequestRequest, opts ...grpc.CallOption) (*pb.CreateRequestResponse, error) {
				return &pb.CreateRequestResponse{}, nil
			},
		},
		{
			name:   "should return bad request #1",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			order:  models.Request{CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022"},
			status: 400,
		},
		{
			name:   "should return bad request #2",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			order:  models.Request{Description: "motor", CustomerId: 8983, DeliveryDate: "10/01/2022"},
			status: 400,
		},
		{
			name:   "should return a general error when server failed to create an request",
			error:  map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			order:  models.Request{Description: "motor", CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022"},
			status: 500,
			server: func(ctx context.Context, in *pb.CreateRequestRequest, opts ...grpc.CallOption) (*pb.CreateRequestResponse, error) {
				return nil, errors.New("some request grpc error")
			},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.order)
		client := &MockRequestServiceClient{}
		client.CreateRequestFunc = tt.server
		CreateRequest(c, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code > 201 {
			b, _ := io.ReadAll(w.Body)
			var errorValue = map[string]string{}
			if err := json.Unmarshal(b, &errorValue); err != nil {
				t.Error(tt.name, "invalid error type", b)
			}

			if w.Code != 201 && !reflect.DeepEqual(errorValue, tt.error) {
				t.Error(tt.name, ",error doesn't match", b)
			}
		}
	}
}

func TestGetRequest(t *testing.T) {

}
