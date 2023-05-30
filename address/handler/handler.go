package handler

import (
	"address-go/helper/response"
	"address-go/model"
	"address-go/service"
	"fmt"
	"net/http"
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
	userIDString, ok := ctx.GetQuery("user_id")
	fmt.Println(ok)
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if userID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("userID must be positive number"))
		return
	}

	res, err := h.svc.Get(userID)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := model.AddressRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if req.UserID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("user_id must be positive number"))
		return
	}

	if req.City == "" {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("city should not be empty"))
		return
	}

	if req.State == "" {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("state should not be empty"))
		return
	}

	if req.Street == "" {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("street should not be empty"))
		return
	}

	if req.Country == "" {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("country should not be empty"))
		return
	}

	if req.Zipcode == "" {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("zipcode should not be empty"))
		return
	}

	res, err := h.svc.Create(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Delete(ctx *gin.Context) {
	addressIDString, ok := ctx.GetQuery("address_id")
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param cart_id should not be empty"))
		return
	}

	userIDString, ok := ctx.GetQuery("user_id")
	fmt.Println(ok)
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	addressID, err := strconv.Atoi(addressIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if userID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("address_id must be positive number"))
		return
	}
	
	if userID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("user_id must be positive number"))
		return
	}

	if err = h.svc.Delete(userID, addressID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}
