package disposition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/disposition/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/bdarge/sb-api-gateway/pkg/utils"
	"google.golang.org/grpc"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockRequestServiceClient struct {
	CreateDispositionFunc func(ctx context.Context, in *pb.CreateDispositionRequest, opts ...grpc.CallOption) (*pb.CreateDispositionResponse, error)
	GetDispositionFunc    func(ctx context.Context, in *pb.GetDispositionRequest, opts ...grpc.CallOption) (*pb.GetDispositionResponse, error)
}

func (m MockRequestServiceClient) CreateDisposition(ctx context.Context, in *pb.CreateDispositionRequest, opts ...grpc.CallOption) (*pb.CreateDispositionResponse, error) {
	return m.CreateDispositionFunc(ctx, in, opts...)
}

func (m MockRequestServiceClient) GetDisposition(ctx context.Context, in *pb.GetDispositionRequest, opts ...grpc.CallOption) (*pb.GetDispositionResponse, error) {
	return m.GetDispositionFunc(ctx, in, opts...)
}

func TestCreateDisposition(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		status       int
		order        models.Disposition
		server       func(ctx context.Context, in *pb.CreateDispositionRequest, opts ...grpc.CallOption) (*pb.CreateDispositionResponse, error)
	}{
		{
			name:   "should create a disposition",
			error:  nil,
			order:  models.Disposition{Description: "motor", CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022", RequestType: "order"},
			status: 201,
			server: func(ctx context.Context, in *pb.CreateDispositionRequest, opts ...grpc.CallOption) (*pb.CreateDispositionResponse, error) {
				return &pb.CreateDispositionResponse{}, nil
			},
		},
		{
			name:   "should return bad disposition #1",
			error:  map[string][]models.ErrorMsg{"errors": {{Field: "Description", Message: "This field is required"}}},
			order:  models.Disposition{CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022", RequestType: "order"},
			status: 400,
		},
		{
			name:   "should return bad disposition #2",
			error:  map[string][]models.ErrorMsg{"errors": {{Field: "CreatedBy", Message: "This field is required"}}},
			order:  models.Disposition{Description: "motor", CustomerId: 8983, DeliveryDate: "10/01/2022", RequestType: "order"},
			status: 400,
		},
		{
			name:         "should return a general error when server failed to create a disposition",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			order:        models.Disposition{Description: "motor", CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022", RequestType: "order"},
			status:       500,
			server: func(ctx context.Context, in *pb.CreateDispositionRequest, opts ...grpc.CallOption) (*pb.CreateDispositionResponse, error) {
				return nil, errors.New("some disposition grpc error")
			},
		},
		{
			name:   "should return bad disposition #3",
			error:  map[string][]models.ErrorMsg{"errors": {{Field: "RequestType", Message: "Should be one of the following: 'order', or 'quote'"}}},
			order:  models.Disposition{Description: "motor", CreatedBy: 12341, CustomerId: 8983, DeliveryDate: "10/01/2022", RequestType: "cat"},
			status: 400,
			server: nil,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.order)
		client := &MockRequestServiceClient{}
		client.CreateDispositionFunc = tt.server
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

}
