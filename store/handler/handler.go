package handler

import (
	"fmt"
	"net/http"
	"store-go/helper/response"
	"store-go/model"
	"store-go/service"
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
	req := []model.StoreRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	for _, v := range req {
		if v.AddressID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("address_id must be positive number"))
			return
		}

		if v.Description == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("description should not be empty"))
			return
		}

		if v.ImageURL == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("image_url should not be empty"))
			return
		}

		if v.Name == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("name should not be empty"))
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
	storeIDString, ok := ctx.GetQuery("store_id")
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param store_id should not be empty"))
		return
	}

	storeID, err := strconv.Atoi(storeIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.svc.Delete(storeID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}
