package user

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/user/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Security BearerAuth
// @Router /users/profile [get]
func (c *Controller) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	user, err := c.service.GetUserByID(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Profile retrieved successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Profile update details"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /users/profile [put]
func (c *Controller) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := c.service.UpdateUser(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Profile updated successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

// GetUser godoc
// @Summary Get user by ID (admin only)
// @Description Get a specific user's information (admin access required)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetUserByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("User retrieved successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary Delete user (admin only)
// @Description Delete a user account (admin access required)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteUser(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Description Get a paginated list of all users (admin access required)
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /users [get]
func (c *Controller) ListUsers(ctx *gin.Context) {
	page := int64(1)
	limit := int64(10)

	if p := ctx.Query("page"); p != "" {
		if parsed, err := strconv.ParseInt(p, 10, 64); err == nil {
			page = parsed
		}
	}

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}

	skip := (page - 1) * limit
	users, err := c.service.ListUsers(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	total := int64(len(users))
	pagination := utils.CalculatePagination(page, limit, total)
	resp := utils.NewSuccessResponseWithPagination("Users retrieved successfully", users, pagination)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) CreateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	profile, err := c.service.CreateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewCreatedResponse("Financial profile created successfully", profile)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) GetFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	profile, err := c.service.GetFinancialProfile(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile retrieved successfully", profile)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UpdateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	profile, err := c.service.UpdateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile updated successfully", profile)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) DeleteFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	err := c.service.DeleteFinancialProfile(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// DisableUser godoc
// @Summary Disable user (admin only)
// @Description Disable a user account (admin access required)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User disabled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /users/{id}/disable [post]
func (c *Controller) DisableUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	err := c.service.DisableUser(ctx, userID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User disabled successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// EnableUser godoc
// @Summary Enable user (admin only)
// @Description Enable a user account (admin access required)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User enabled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /users/{id}/enable [post]
func (c *Controller) EnableUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	err := c.service.EnableUser(ctx, userID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User enabled successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
