package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bdarge/sb-api-gateway/pkg/auth/pb"
	"github.com/bdarge/sb-api-gateway/pkg/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type MockAuthServiceClient struct {
	LoginFunc    func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error)
	RegisterFunc func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error)
	ValidateFunc func(ctx context.Context, in *pb.ValidateRequest, opts ...grpc.CallOption) (*pb.ValidateResponse, error)
}

func (M MockAuthServiceClient) Register(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
	return M.RegisterFunc(ctx, in, opts...)
}

func (M MockAuthServiceClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	return M.LoginFunc(ctx, in, opts...)
}

func (M MockAuthServiceClient) Validate(ctx context.Context, in *pb.ValidateRequest, opts ...grpc.CallOption) (*pb.ValidateResponse, error) {
	return M.ValidateFunc(ctx, in, opts...)
}

func TestLogin(t *testing.T) {
	// tests cases
	tests := []struct {
		name      string
		error     string
		status    int
		body      models.Account
		dependent func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error)
		result    map[string]string
	}{
		{
			name:   "remote server error",
			error:  "some error",
			status: 500,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
				return nil, errors.New("some error")
			},
			result: nil,
		},
		{
			name:   "invalid post data",
			error:  "some error",
			status: 400,
			body:   models.Account{Email: "", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "invalid email data",
			error:  "some error",
			status: 400,
			body:   models.Account{Email: "etwyte", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "happy path",
			error:  "",
			status: 200,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
				return &pb.LoginResponse{Token: "some_value"}, nil
			},
			result: map[string]string{"token": "some_value"},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx := mockPostTest(w, tt.body)
		client := new(MockAuthServiceClient)
		client.LoginFunc = tt.dependent
		Login(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}
		if tt.status == 200 {
			b, _ := io.ReadAll(w.Body)
			var jsonMap pb.LoginResponse
			json.Unmarshal([]byte(b), &jsonMap)

			if jsonMap.Token != tt.result["token"] {
				t.Error(tt.name, "Token is missing", w.Code, string(b))
			}
		}
	}
}

func TestRegister(t *testing.T) {
	// tests cases
	tests := []struct {
		name      string
		error     map[string]string
		status    int
		body      models.Account
		dependent func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error)
		result    map[string]string
	}{
		{
			name:   "remote server error",
			error:  map[string]string{"error": "VALIDATEERR-2", "message": "Error happened at the server, please check later."},
			status: 500,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
				return nil, errors.New("some error")
			},
			result: nil,
		},
		{
			name:   "invalid post data",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			status: 400,
			body:   models.Account{Email: "", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "invalid email data",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			status: 400,
			body:   models.Account{Email: "etwyte", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "happy path",
			error:  nil,
			status: 201,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
				return &pb.RegisterResponse{Status: http.StatusCreated}, nil
			},
			result: map[string]string{"token": "some_value"},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := mockPostTest(w, tt.body)
		client := &MockAuthServiceClient{}
		client.RegisterFunc = tt.dependent

		Register(c, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
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

func mockPostTest(w *httptest.ResponseRecorder, content interface{}) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ctx.Request.Method = "POST"

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

	return ctx
}
