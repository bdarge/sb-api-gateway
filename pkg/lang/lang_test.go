package lang

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	. "github.com/bdarge/api-gateway/out/lang"
	"github.com/bdarge/api-gateway/out/model"
	"github.com/bdarge/api-gateway/pkg/models"
	"github.com/bdarge/api-gateway/pkg/utils"
	"google.golang.org/grpc"
)


type MockRequestServiceClient struct{
	GetLangFunc func(ctx context.Context, in *LangGetRequest, opts ...grpc.CallOption) (*LangGetResponse, error)
}

func (m MockRequestServiceClient) GetLang(ctx context.Context, in *LangGetRequest, opts ...grpc.CallOption) (*LangGetResponse, error) {
	return m.GetLangFunc(ctx, in, opts...)
}

func TestGetLang(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		er           string
		error        map[string][]models.ErrorMsg
		generalError map[string]string
		status       int
		data 				 *LangGetResponse
		responseData models.Langs
	}{
		{
			name:  "should return all lang",
			er: "",
			error: nil,
			status: 200,
			data: &LangGetResponse{
				Data: []*model.LangData{
					{
					Id: 1,
					Language: "en",
					Currency: "usd",
				},
			}},
			responseData: models.Langs{
				Data: []models.Lang{
					{
					ID: 1,
					Language:  "en",
					Currency: "usd",
				},
			}},
		},
		{
			name:"should return an error:",
			generalError: map[string]string{"error": "ACTIONERR-1", "message": "An error happened, please check later."},
			status: 500,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		client := &MockRequestServiceClient{}
		client.GetLangFunc = func(ctx context.Context, in *LangGetRequest, opts ...grpc.CallOption) (*LangGetResponse, error) {
			if tt.status > 200 {
				return nil, errors.New("some backend service grpc error")
			}
			return tt.data, nil
		}
		ctx := utils.GetTestGinContext(w)

		utils.MockGetTest(ctx, nil, url.Values{})

		// act
		GetLang(ctx, client)

		if w.Code != tt.status {
			b, _ := io.ReadAll(w.Body)
			t.Error(tt.name, w.Code, string(b))
			continue
		}

		if w.Code == 200 {
			b, _ := io.ReadAll(w.Body)
			response := &models.Langs{}
			err := json.Unmarshal(b, response)
			if err != nil {
				t.Error(tt.name, "test error:", err)
				continue
			}
			if !reflect.DeepEqual(*response, tt.responseData) {
				t.Error(tt.name, "data doesn't match actual:", string(b), "expected:", tt.responseData)
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
