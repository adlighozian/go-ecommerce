package handler

import "github.com/gin-gonic/gin"

type UserSettingHandlerI interface {
	UpdateByUserID(ctx *gin.Context)
}
