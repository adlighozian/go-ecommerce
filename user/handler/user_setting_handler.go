package handler

import (
	"net/http"
	"strconv"
	"user-go/helper/response"
	"user-go/model"
	"user-go/service"

	"github.com/gin-gonic/gin"
)

type UserSettingHandler struct {
	svc service.UserSettingServiceI
}

func NewUserSettingHandler(svc service.UserSettingServiceI) UserSettingHandlerI {
	h := new(UserSettingHandler)
	h.svc = svc
	return h
}

func (h *UserSettingHandler) UpdateByUserID(ctx *gin.Context) {
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

	settingReq := new(model.SettingReq)
	if bindErr := ctx.ShouldBindJSON(&settingReq); bindErr != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", bindErr.Error())
		return
	}

	user, errSvc := h.svc.UpdateByUserID(uint(userID), settingReq)
	if errSvc != nil {
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"user": user,
	})
}
