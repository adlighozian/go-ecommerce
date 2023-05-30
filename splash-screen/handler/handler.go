package handler

import (
	"fmt"
	"net/http"
	"splash-screen-go/helper/response"
	"splash-screen-go/model"
	"splash-screen-go/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service.Servicer
}

func NewHandler(svc service.Servicer) Handlerer {
	return &handler{
		svc: svc,
	}
}

func (h *handler) Get(ctx *gin.Context) {
	res, err := h.svc.Get()
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := []model.SplashScreenRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	for _, v := range req {
		if v.ImageURL == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("image_url should not be empty"))
			return
		}
	}

	res, err := h.svc.Create(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Delete(ctx *gin.Context) {
	splashScreenIDString, ok := ctx.GetQuery("splash_id")
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param splash_id should not be empty"))
		return
	}

	splashScreenID, err := strconv.Atoi(splashScreenIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.svc.Delete(splashScreenID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}
