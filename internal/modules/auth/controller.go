package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/auth/dto"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := c.service.Register(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewCreatedResponse("User registered successfully", dto.RegisterResponse{User: user})
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	authResp, err := c.service.Login(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, err.Error())
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	resp := utils.NewSuccessResponse("Login successful", authResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	accessToken, err := c.service.RefreshAccessToken(ctx, userID.(string), req.RefreshToken)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, err.Error())
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	resp := utils.NewSuccessResponse("Token refreshed successfully", dto.RefreshTokenResponse{AccessToken: accessToken})
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	err := c.service.Logout(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Logged out successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
