package handler

import (
	"auth-go/helper/authjwt"
	"auth-go/helper/response"
	"auth-go/model"
	"auth-go/service"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	svc          service.AuthServiceI
	jwtSecretKey string
	gauth        *oauth2.Config
}

func NewAuthHandler(svc service.AuthServiceI, jwtSecretKey string, gauth *oauth2.Config) AuthHandlerI {
	h := new(AuthHandler)
	h.svc = svc
	h.jwtSecretKey = jwtSecretKey
	h.gauth = gauth
	return h
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	registerReq := new(model.RegisterReq)
	if bindErr := ctx.ShouldBindJSON(&registerReq); bindErr != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", bindErr.Error())
		return
	}

	user, errSvc := h.svc.Create(registerReq)
	if errSvc != nil {
		if errors.Is(errSvc, errors.New("duplicated key not allowed")) {
			response.NewJSONResErr(ctx, http.StatusConflict, "", "email already existed")
			return
		}
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	response.NewJSONRes(ctx, http.StatusCreated, "", user.ID)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	loginReq := new(model.LoginReq)
	if bindErr := ctx.ShouldBindJSON(&loginReq); bindErr != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", bindErr.Error())
		return
	}

	user, errSvc := h.svc.GetByEmail(loginReq)
	if errSvc != nil {
		if errors.Is(errSvc, errors.New("incorrect password")) {
			response.NewJSONResErr(ctx, http.StatusBadRequest, "", errSvc.Error())
			return
		}
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	accessTokenDur := 24 * time.Hour
	accessToken, errGenAccToken := authjwt.GenerateToken(
		h.jwtSecretKey, accessTokenDur,
		user.ID, user.Role,
	)
	if errGenAccToken != nil {
		_ = ctx.Error(errGenAccToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGenAccToken.Error())
		return
	}

	refreshTokenDur := 7 * 24 * time.Hour
	refreshToken, errGenRefToken := authjwt.GenerateToken(
		h.jwtSecretKey, refreshTokenDur,
		user.ID, user.Role,
	)
	if errGenRefToken != nil {
		_ = ctx.Error(errGenRefToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGenRefToken.Error())
		return
	}

	ctx.SetCookie("user_id", strconv.FormatUint(uint64(user.ID), 10), 0, "", "", true, true)
	ctx.SetCookie("user_role", user.Role, 0, "", "", true, true)

	ctx.SetCookie("access_token", accessToken, 0, "", "", true, true)

	// Store the refresh token in Redis
	dataByte, _ := json.Marshal(model.RefreshToken{
		UserID:   user.ID,
		UserRole: user.Role,
	})
	errSetRefToken := h.svc.SetRefreshToken(refreshToken, dataByte, refreshTokenDur)
	if errSetRefToken != nil {
		_ = ctx.Error(errSetRefToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSetRefToken.Error())
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "", "", true, true)

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	// Get the refresh token from the request cookies
	cookie, errCookie := ctx.Cookie("refresh_token")
	if errCookie != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "Refresh token required", errCookie.Error())
		return
	}

	refreshToken, errSvc := h.svc.GetByRefreshToken(cookie)
	if errSvc != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "Invalid refresh token", errSvc.Error())
		return
	}

	accessTokenDur := 24 * time.Hour
	accessToken, errGenAccToken := authjwt.GenerateToken(
		h.jwtSecretKey, accessTokenDur,
		refreshToken.UserID, refreshToken.UserRole,
	)
	if errGenAccToken != nil {
		_ = ctx.Error(errGenAccToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGenAccToken.Error())
		return
	}

	ctx.SetCookie("access_token", accessToken, 0, "", "", true, true)

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) GoogleLogin(ctx *gin.Context) {
	url := h.gauth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(ctx *gin.Context) {
	code, exist := ctx.GetQuery("code")
	if !exist {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "", "code required")
		return
	}

	token, errExc := h.gauth.Exchange(context.Background(), code)
	if errExc != nil {
		response.NewJSONResErr(ctx, http.StatusBadRequest, "Failed to exchange token", errExc.Error())
		return
	}

	userReq := new(model.UserReq)

	req, errReq := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil,
	)
	if errReq != nil {
		_ = ctx.Error(errReq)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errReq.Error())
		return
	}

	client := h.gauth.Client(context.Background(), token)
	resp, errResp := client.Do(req)
	if errResp != nil {
		_ = ctx.Error(errResp)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errResp.Error())
		return
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		_ = ctx.Error(errRead)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errRead.Error())
		return
	}

	var data map[string]interface{}
	errJSONUn := json.Unmarshal(body, &data)
	if errJSONUn != nil {
		_ = ctx.Error(errJSONUn)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errJSONUn.Error())
		return
	}

	userReq.Username, _ = data["name"].(string)
	userReq.Email, _ = data["email"].(string)
	userReq.FullName, _ = data["name"].(string)
	userReq.ImageURL, _ = data["picture"].(string)

	// user := &model.UserReq{
	// 	Username: data["name"].(string),
	// 	Email:    data["email"].(string),
	// 	FullName: data["name"].(string),
	// 	ImageURL: data["picture"].(string),
	// }
	// userReq, errGetUser := h.getUser(client)
	// if errGetUser != nil {
	// 	_ = ctx.Error(errGetUser)
	// 	response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGetUser.Error())
	// 	return
	// }

	user, errSvc := h.svc.FirstOrCreate(userReq)
	if errSvc != nil {
		_ = ctx.Error(errSvc)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSvc.Error())
		return
	}

	accessTokenDur := 24 * time.Hour
	accessToken, errGenAccToken := authjwt.GenerateToken(
		h.jwtSecretKey, accessTokenDur,
		user.ID, user.Role,
	)
	if errGenAccToken != nil {
		_ = ctx.Error(errGenAccToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGenAccToken.Error())
		return
	}

	refreshTokenDur := 7 * 24 * time.Hour
	refreshToken, errGenRefToken := authjwt.GenerateToken(
		h.jwtSecretKey, refreshTokenDur,
		user.ID, user.Role,
	)
	if errGenRefToken != nil {
		_ = ctx.Error(errGenRefToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errGenRefToken.Error())
		return
	}

	ctx.SetCookie("user_id", strconv.FormatUint(uint64(user.ID), 10), 0, "", "", true, true)
	ctx.SetCookie("user_role", user.Role, 0, "", "", true, true)

	ctx.SetCookie("access_token", accessToken, 0, "", "", true, true)

	// Store the refresh token in Redis
	dataByte, _ := json.Marshal(model.RefreshToken{
		UserID:   user.ID,
		UserRole: user.Role,
	})
	errSetRefToken := h.svc.SetRefreshToken(refreshToken, dataByte, refreshTokenDur)
	if errSetRefToken != nil {
		_ = ctx.Error(errSetRefToken)
		response.NewJSONResErr(ctx, http.StatusInternalServerError, "", errSetRefToken.Error())
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 0, "", "", true, true)

	response.NewJSONRes(ctx, http.StatusOK, "", map[string]any{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// func (h *AuthHandler) getUser(client *http.Client) (*model.UserReq, error) {
// 	req, errReq := http.NewRequestWithContext(
// 		context.Background(),
// 		http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil,
// 	)
// 	if errReq != nil {
// 		return nil, errReq
// 	}

// 	resp, errResp := client.Do(req)
// 	if errResp != nil {
// 		return nil, errResp
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data map[string]interface{}
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &model.UserReq{
// 		Username: data["name"].(string),
// 		Email:    data["email"].(string),
// 		FullName: data["name"].(string),
// 		ImageURL: data["picture"].(string),
// 	}

// 	return user, nil
// }
