package handler

import (
	"fmt"
	"net/http"
	"payment-go/helper/response"
	"payment-go/model"
	"payment-go/service"
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
		if v.Name == "" {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("payment method's name should not be empty"))
			return
		}

		if v.PaymentGatewayID == "" {
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
	req := model.PaymentLogRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if req.UserID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("user_id should be positive number"))
		return
	}

	if req.OrderID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("order_id should be positive number"))
		return
	}

	if req.PaymentMethodID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("payment_method_id should not be empty"))
		return
	}

	res, err := h.svc.CreatePaymentLog(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) DeletePaymentMethod(ctx *gin.Context) {
	paymentMethodIDString, ok := ctx.GetQuery("payment_method_id")
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param payment_method_id should not be empty"))
		return
	}

	paymentMethodID, err := strconv.Atoi(paymentMethodIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.svc.DeletePaymentMethod(paymentMethodID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}
