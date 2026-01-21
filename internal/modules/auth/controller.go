package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/HasanNugroho/coin-be/internal/modules/auth/dto"
	"github.com/HasanNugroho/coin-be/pkg/errors"
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
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	user, err := c.service.Register(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, dto.RegisterResponse{User: user})
}

func (c *Controller) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	authResp, err := c.service.Login(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, authResp)
}

func (c *Controller) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	accessToken, err := c.service.RefreshAccessToken(ctx, userID.(string), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.RefreshTokenResponse{AccessToken: accessToken})
}

func (c *Controller) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	err := c.service.Logout(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
