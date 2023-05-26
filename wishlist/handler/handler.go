package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"wishlist-go/helper/response"
	"wishlist-go/model"
	"wishlist-go/service"
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
	userIDString, ok := ctx.GetQuery("user_ID")	
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

func (h *handler) GetDetail(ctx *gin.Context) {
	userIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	wishlistIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param wishlist_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	wishlistID, err := strconv.Atoi(wishlistIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}
		
	if userID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("userID must be positive number"))
		return
	}

	if wishlistID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("wishlistID must be positive number"))
		return
	}

	res, err := h.svc.GetDetail(userID, wishlistID)
	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, res)
}

func (h *handler) Create(ctx *gin.Context) {
	req := []model.WishlistRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	for _, v := range req {
		if v.ProductID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("producID must be positive number"))
			return
		}

		if v.UserID <= 0 {
			response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("producID must be positive number"))
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
	userIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param user_id should not be empty"))
		return
	}

	wishlistIDString, ok := ctx.GetQuery("user_ID")	
	if !ok {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("query param wishlist_id should not be empty"))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	wishlistID, err := strconv.Atoi(wishlistIDString)
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	if userID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("userID must be positive number"))
		return
	}

	if wishlistID <= 0 {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Errorf("wishlistID must be positive number"))
		return
	}

	if err = h.svc.Delete(userID, wishlistID); err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, err)
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, nil)
}