package handler

import (
	"api-gateway-go/helper/response"
	"api-gateway-go/helper/timeout"
	"api-gateway-go/model"
	"api-gateway-go/service"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortenHandler struct {
	svc service.ShortenServiceI
}

func NewShortenHandler(svc service.ShortenServiceI) ShortenHandlerI {
	h := new(ShortenHandler)
	h.svc = svc
	return h
}

func (h *ShortenHandler) Get(ctx *gin.Context) {
	var url string
	urlCtx, _ := ctx.Get("url")
	url, _ = urlCtx.(string)

	var userID string
	userIDCtx, _ := ctx.Get("userID")
	userID, _ = userIDCtx.(string)

	timeoutCtx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	req, errReq := http.NewRequestWithContext(timeoutCtx, ctx.Request.Method, url, ctx.Request.Body)
	if errReq != nil {
		_ = ctx.Error(errReq)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errReq.Error())
		return
	}
	req.Header.Add("user-id", userID)

	resp, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		_ = ctx.Error(errResp)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errResp.Error())
		return
	}
	defer resp.Body.Close()

	// Copy the response headers to the gateway response
	h.copyHeaders(ctx.Writer.Header(), resp.Header)

	// Set the status code
	ctx.Writer.WriteHeader(resp.StatusCode)

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		_ = ctx.Error(errRead)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errRead.Error())
		return
	}

	var jsonRes response.JSONRes
	if errJSONUn := json.Unmarshal(body, &jsonRes); errJSONUn != nil {
		_ = ctx.Error(errJSONUn)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errJSONUn.Error())
		return
	}

	ctx.JSON(http.StatusOK, jsonRes)
}

func (h *ShortenHandler) copyHeaders(dst http.Header, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func (h *ShortenHandler) Shorten(ctx *gin.Context) {
	shortenReq := new(model.ShortenReq)
	if errJSON := ctx.ShouldBindJSON(&shortenReq); errJSON != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", errJSON.Error())
		return
	}

	apiManagement, errSvc := h.svc.Create(shortenReq)
	if errSvc != nil {
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	response.NewJSONRes(ctx, http.StatusOK, "", apiManagement)
}
