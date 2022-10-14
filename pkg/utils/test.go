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

func MockPostTest(w *httptest.ResponseRecorder, content interface{}) *gin.Context {
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
