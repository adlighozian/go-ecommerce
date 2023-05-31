package handler

import (
	"net/http"
	"user-go/helper/response"

	"github.com/gin-gonic/gin"
)

type PingGinHandler struct{}

func NewPingGinHandler() PingGinHandlerI {
	return new(PingGinHandler)
}

type PingGinHandlerI interface {
	Ping(c *gin.Context)
}

func (h *PingGinHandler) Ping(c *gin.Context) {
	response.NewJSONRes(c, http.StatusOK, "", "pong")
}
