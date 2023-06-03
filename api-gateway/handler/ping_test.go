package handler_test

import (
	"api-gateway-go/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingGinHandler_Ping(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		// TODO: Add test cases.
		{name: "success", wantStatus: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			h := handler.NewPingGinHandler()

			req, errReq := http.NewRequest(http.MethodGet, "/ping", nil)
			assert.NoError(t, errReq)

			res := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(res)
			c.Request = req

			h.Ping(c)

			t.Log(res.Code, res.Body.String())
			t.Log("------------------------------------\n")

			assert.Equal(t, tt.wantStatus, res.Code, "status code got: %v, want %v", res.Code, tt.wantStatus)
		})
	}
}
