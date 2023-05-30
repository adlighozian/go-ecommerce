package handler

import (
	"api-gateway-go/helper/response"
	"api-gateway-go/helper/timeout"
	"api-gateway-go/model"
	"api-gateway-go/service"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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

	reqBody := ctx.Request.Body
	req, errReq := http.NewRequestWithContext(timeoutCtx, ctx.Request.Method, url, reqBody)
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

	respBody, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		_ = ctx.Error(errRead)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errRead.Error())
		return
	}

	// if google oauth2 is called, return to []byte / html
	if strings.Contains(url, "/google/") {
		ctx.Data(http.StatusOK, "text/html; charset=UTF-8", respBody)
		return
	}

	var jsonRes response.JSONRes
	if errJSONUn := json.Unmarshal(respBody, &jsonRes); errJSONUn != nil {
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

// credits: https://github.com/stripe/stripe-go/blob/cb7a4cc7ba3ad39fa36d6bfafc9e89b8e6350b05/stripe.go#L318
// nopReadCloser's sole purpose is to give us a way to turn an `io.Reader` into
// an `io.ReadCloser` by adding a no-op implementation of the `Closer`
// interface. We need this because `http.Request`'s `Body` takes an
// `io.ReadCloser` instead of a `io.Reader`.
// type nopReadCloser struct {
// 	io.Reader
// }

// func (nopReadCloser) Close() error { return nil }
