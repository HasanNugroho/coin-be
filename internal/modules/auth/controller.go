package auth

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/auth/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
	userSrv *user.Service
}

func NewController(service *Service, userSrv *user.Service) *Controller {
	return &Controller{service: service, userSrv: userSrv}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /v1/auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/auth/login [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token details"
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/auth/refresh-token [post]
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

// Logout godoc
// @Summary User logout
// @Description Invalidate the user's refresh token and end the session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logged out successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/auth/logout [post]
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

// GetMe godoc
// @Summary Get current user profile
// @Description Retrieve the authenticated user's complete profile information including user and profile data
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User profile retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/auth/me [get]
func (c *Controller) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	user, err := c.userSrv.GetUserProfile(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("fetch user profile successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

// ValidateToken godoc
// @Summary Validate token
// @Description Check if the provided token is valid
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} map[string]interface{} "Token is invalid"
// @Security BearerAuth
// @Router /v1/auth/validate [get]
func (c *Controller) ValidateToken(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "invalid token")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	resp := utils.NewSuccessResponse("token is valid", map[string]interface{}{
		"valid":   true,
		"user_id": userID,
	})
	ctx.JSON(http.StatusOK, resp)
}
