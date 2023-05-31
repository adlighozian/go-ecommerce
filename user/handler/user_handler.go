package handler

import (
	"net/http"
	"strconv"
	"user-go/helper/response"
	"user-go/model"
	"user-go/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserServiceI
}

func NewUserHandler(svc service.UserServiceI) UserHandlerI {
	h := new(UserHandler)
	h.svc = svc
	return h
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	userIDStr := ctx.GetHeader("user-id")
	if userIDStr == "" {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", "missing user-id header")
		return
	}

	userID, errParse := strconv.ParseUint(userIDStr, 10, 32)
	if userIDStr == "" {
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errParse.Error())
		return
	}

	user, errSvc := h.svc.GetByID(uint(userID))
	if errSvc != nil {
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"user": user,
	})
}

func (h *UserHandler) UpdateByID(ctx *gin.Context) {
	userIDStr := ctx.GetHeader("user-id")
	if userIDStr == "" {
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", "missing user-id header")
		return
	}

	userID, errParse := strconv.ParseUint(userIDStr, 10, 32)
	if userIDStr == "" {
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errParse.Error())
		return
	}

	profileReq := new(model.ProfileReq)
	if bindErr := ctx.ShouldBindJSON(&profileReq); bindErr != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", bindErr.Error())
		return
	}

	user, errSvc := h.svc.UpdateByID(uint(userID), profileReq)
	if errSvc != nil {
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"user": user,
	})
}
