package handler_test

import (
	"api-gateway-go/handler"
	"api-gateway-go/mocks"
	"api-gateway-go/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenHandler_Shorten(t *testing.T) {
	isAvailable, needBypass := true, false
	type args struct {
		req *model.ShortenReq
	}
	tests := []struct {
		name       string
		args       args
		ucRes      *model.APIManagement
		ucErr      error
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				req: &model.ShortenReq{
					APIName:     "wislists",
					ServiceName: "wislists",
					EndpointURL: "http://localhost:5013/wislists",
					IsAvailable: &isAvailable,
					NeedBypass:  &needBypass,
				},
			},
			ucRes: &model.APIManagement{
				ID:                1,
				APIName:           "wislists",
				ServiceName:       "wislists",
				EndpointURL:       "http://localhost:5013/wislists",
				HashedEndpointURL: "CUbICN8VgNk",
				IsAvailable:       true,
				NeedBypass:        false,
				CreatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
				UpdatedAt:         time.Date(2023, 6, 6, 10, 20, 0, 0, time.Local),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "failed svc error",
			args: args{
				req: &model.ShortenReq{
					APIName:     "wislists",
					ServiceName: "wislists",
					EndpointURL: "http://localhost:5013/wislists",
					IsAvailable: &isAvailable,
					NeedBypass:  &needBypass,
				},
			},
			ucErr:      errors.New("svc error"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := mocks.NewShortenServiceI(t)

			h := handler.NewShortenHandler(mockUC)
			if tt.ucRes != nil || tt.ucErr != nil {
				mockUC.On("Create", mock.Anything).Return(tt.ucRes, tt.ucErr)
			}

			jsonBytes, _ := json.Marshal(tt.args.req)
			reqBody := bytes.NewBuffer(jsonBytes)

			req, errReq := http.NewRequest(http.MethodPost, "/shorten", reqBody)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.Shorten(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code got: %v, want %v", res.Code, tt.wantStatus)
		})
	}
}

type brokenReader struct{}

func (br brokenReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestShortenHandler_Get(t *testing.T) {
	type args struct {
		hashedURL string
	}
	tests := []struct {
		name       string
		args       args
		mockHTTP   func(req *http.Request) (*http.Response, error)
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				hashedURL: "hashedURL1",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, `{"message":"success"}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not_found",
			args: args{
				hashedURL: "hashedURL2",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(404, `{"message":"Not Found"}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal_server_error",
			args: args{
				hashedURL: "hashedURL3",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(500, `{"message":"Internal Server Error"}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid_json_response",
			args: args{
				hashedURL: "hashedURL4",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, `{"message":}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			// give on this one lol
			name: "timeout",
			args: args{
				hashedURL: "hashedURL5",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				time.Sleep(30 * time.Second)
				return nil, context.DeadlineExceeded
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "http_request_error",
			args: args{
				hashedURL: "hashedURL6",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("http request error")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "read_body_error",
			args: args{
				hashedURL: "hashedURL7",
			},
			mockHTTP: func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, `{"message":"success"}`)
				resp.Header.Set("Content-Type", "application/json")
				resp.Body = io.NopCloser(brokenReader{})
				return resp, nil
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			// Activate the httpmock
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			url := "http://test.com/" + tt.args.hashedURL

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			c.Request = req

			// Setting the URL as a context value
			c.Set("url", url)

			mockUC := mocks.NewShortenServiceI(t)

			h := handler.NewShortenHandler(mockUC)

			// Set up httpmock to intercept the HTTP request and return a successful response
			httpmock.RegisterResponder(http.MethodGet, url, tt.mockHTTP)

			h.Get(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code got: %v, want %v", res.Code, tt.wantStatus)
		})
	}
}
