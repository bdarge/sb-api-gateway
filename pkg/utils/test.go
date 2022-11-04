package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func MockPostTest(w *httptest.ResponseRecorder, content interface{}) *gin.Context {
	ctx := GetTestGinContext(w)
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

func MockGetTest(ctx *gin.Context, params gin.Params, u url.Values) {
	ctx.Request.Method = "GET"
	ctx.Request.Header.Set("Content-Type", "application/json")

	u.Add("skip", "5")
	u.Add("limit", "10")

	// set path params
	ctx.Params = params

	//set query params
	ctx.Request.URL.RawQuery = u.Encode()
}
