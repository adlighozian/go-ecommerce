package handler

import (
	"fmt"
	"net/http"
	"payment-go/helper/response"
	"payment-go/model"
	"payment-go/service"

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

func (h *handler) GetPaymentMethod(ctx *gin.Context) {
	res, err := h.svc.GetPaymentMethod()
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) CreatePaymentMethod(ctx *gin.Context) {
	req := []model.PaymentMethodRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	for _, v := range req {
		if v.Name != "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("payment method's name should not be empty"))
			return
		}

		if v.PaymentGatewayID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("payment_method_id should not be empty"))
			return
		}
	}

	res, err := h.svc.CreatePaymentMethod(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) CreatePaymentLog(ctx *gin.Context) {
	req := []model.PaymentLogsRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	for _, v := range req {
		if v.UserID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("user_id should be positive number"))
			return
		}

		if v.OrderID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("order_id should be positive number"))
			return
		}

		if v.PaymentMethodID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("payment_method_id should not be empty"))
			return
		}
	}

	res, err := h.svc.CreatePaymentLog(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}