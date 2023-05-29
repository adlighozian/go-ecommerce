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

func (h *handler) CheckTransaction(ctx *gin.Context) {
	orderID, ok := ctx.GetQuery("order_id")
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query order_id should not be empty"))
		return
	}

	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("order_id should be number"))
		return
	}

	if orderIDInt <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("order_id should be positive number"))
		return
	}

	res, err := h.svc.CheckTransaction(orderID)
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

	if req.TotalPayment <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("total payment should not be empty"))
		return
	}

	res, err := h.svc.CreatePaymentLog(req)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}