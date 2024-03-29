package auth

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/bdarge/api-gateway/out/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockAuthServiceClient struct {
	LoginFunc        func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	RegisterFunc     func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	ValidateFunc     func(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error)
	RefreshTokenFunc func(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

func (M MockAuthServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	return M.RegisterFunc(ctx, in, opts...)
}

func (M MockAuthServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	return M.LoginFunc(ctx, in, opts...)
}

func (M MockAuthServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) {
	return M.ValidateFunc(ctx, in, opts...)
}

func (M MockAuthServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	return M.RefreshTokenFunc(ctx, in, opts...)
}

func TestLogin(t *testing.T) {
	c, _ := config.LoadConfig("dev")
	// tests cases
	tests := []struct {
		name      string
		error     string
		status    int
		body      models.Account
		dependent func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
		result    map[string]string
	}{
		{
			name:   "remote server error",
			error:  "some error",
			status: 500,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
				return nil, errors.New("some error")
			},
			result: nil,
		},
		{
			name:   "invalid post data",
			error:  "some error",
			status: 400,
			body:   models.Account{Email: "", Password: "some_value"},
			dependent: func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "invalid email data",
			error:  "some error",
			status: 400,
			body:   models.Account{Email: "etwyte", Password: "some_value"},
			dependent: func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "happy path",
			error:  "",
			status: 200,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
				return &LoginResponse{Token: "some_value"}, nil
			},
			result: map[string]string{"token": "some_value"},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx := utils.MockPostTest(w, tt.body)
		client := new(MockAuthServiceClient)
		client.LoginFunc = tt.dependent
		Login(ctx, client, &c)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}
		if tt.status == 200 {
			b, _ := io.ReadAll(w.Body)
			var jsonMap LoginResponse
			json.Unmarshal([]byte(b), &jsonMap)

			if jsonMap.Token != tt.result["token"] {
				t.Error(tt.name, "Token is missing", w.Code, string(b))
			}
		}
	}
}

func TestRegister(t *testing.T) {
	// test cases
	tests := []struct {
		name      string
		error     map[string]string
		status    int
		body      models.Account
		dependent func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
		result    map[string]string
	}{
		{
			name:   "remote server error",
			error:  map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			status: 500,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
				return nil, errors.New("some error")
			},
			result: nil,
		},
		{
			name:   "invalid post data",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			status: 400,
			body:   models.Account{Email: "", Password: "some_value"},
			dependent: func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "invalid email data",
			error:  map[string]string{"error": "VALIDATEERR-1", "message": "Invalid inputs. Please check your inputs"},
			status: 400,
			body:   models.Account{Email: "etwyte", Password: "some_value"},
			dependent: func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
				return nil, nil
			},
			result: nil,
		},
		{
			name:   "happy path",
			error:  nil,
			status: 201,
			body:   models.Account{Email: "fake@fake.com", Password: "some_value"},
			dependent: func(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
				return &RegisterResponse{Status: http.StatusCreated}, nil
			},
			result: map[string]string{"token": "some_value"},
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c := utils.MockPostTest(w, tt.body)
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
