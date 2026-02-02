package user_platform

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service      *Service
	platformRepo *platform.Repository
}

func NewController(s *Service, pr *platform.Repository) *Controller {
	return &Controller{
		service:      s,
		platformRepo: pr,
	}
}

// CreateUserPlatform godoc
// @Summary Create a new user platform
// @Description Create a new user platform for the authenticated user
// @Tags User Platforms
// @Accept json
// @Produce json
// @Param request body dto.CreateUserPlatformRequest true "Platform details"
// @Success 201 {object} map[string]interface{} "User platform created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-platforms [post]
func (c *Controller) CreateUserPlatform(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateUserPlatformRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := utils.ValidateRequest(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatform, err := c.service.CreateUserPlatform(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResp := c.mapToResponse(ctx, userPlatform)
	resp := utils.NewSuccessResponse("User platform created successfully", userPlatformResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetUserPlatform godoc
// @Summary Get user platform by ID
// @Description Get a specific user platform by ID
// @Tags User Platforms
// @Produce json
// @Param id path string true "User Platform ID"
// @Success 200 {object} map[string]interface{} "User platform retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/user-platforms/{id} [get]
func (c *Controller) GetUserPlatform(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	userPlatform, err := c.service.GetUserPlatformByID(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResp := c.mapToResponse(ctx, userPlatform)
	resp := utils.NewSuccessResponse("User platform retrieved successfully", userPlatformResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListUserPlatforms godoc
// @Summary List all user platforms
// @Description Get all user platforms for the authenticated user
// @Tags User Platforms
// @Produce json
// @Success 200 {object} map[string]interface{} "User platforms retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-platforms [get]
func (c *Controller) ListUserPlatforms(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userPlatforms, err := c.service.ListUserPlatforms(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResps := c.mapToResponseList(ctx, userPlatforms)
	resp := utils.NewSuccessResponse("User platforms retrieved successfully", userPlatformResps)
	ctx.JSON(http.StatusOK, resp)
}

// ListUserPlatformsDropdown godoc
// @Summary List user platforms for dropdown
// @Description Get all user platforms for dropdown with platform data lookup
// @Tags User Platforms
// @Produce json
// @Success 200 {object} map[string]interface{} "User platforms retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-platforms/dropdown/list [get]
func (c *Controller) ListUserPlatformsDropdown(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userPlatforms, err := c.service.ListUserPlatformsDropdown(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResps := c.mapToDropdownResponseList(ctx, userPlatforms)
	resp := utils.NewSuccessResponse("User platforms retrieved successfully", userPlatformResps)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateUserPlatform godoc
// @Summary Update user platform
// @Description Update a user platform
// @Tags User Platforms
// @Accept json
// @Produce json
// @Param id path string true "User Platform ID"
// @Param request body dto.UpdateUserPlatformRequest true "Update details"
// @Success 200 {object} map[string]interface{} "User platform updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/user-platforms/{id} [put]
func (c *Controller) UpdateUserPlatform(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	var req dto.UpdateUserPlatformRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := utils.ValidateRequest(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatform, err := c.service.UpdateUserPlatform(ctx, userID.(string), id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResp := c.mapToResponse(ctx, userPlatform)
	resp := utils.NewSuccessResponse("User platform updated successfully", userPlatformResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteUserPlatform godoc
// @Summary Delete user platform
// @Description Delete a user platform
// @Tags User Platforms
// @Produce json
// @Param id path string true "User Platform ID"
// @Success 200 {object} map[string]interface{} "User platform deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/user-platforms/{id} [delete]
func (c *Controller) DeleteUserPlatform(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	err := c.service.DeleteUserPlatform(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User platform deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(ctx *gin.Context, userPlatform *UserPlatform) *dto.UserPlatformResponse {
	var platformData *dto.PlatformData
	platform, err := c.platformRepo.GetPlatformByID(ctx, userPlatform.PlatformID)
	if err == nil && platform != nil {
		platformData = &dto.PlatformData{
			ID:       platform.ID.Hex(),
			Name:     platform.Name,
			Type:     platform.Type,
			IsActive: platform.IsActive,
		}
	}

	return &dto.UserPlatformResponse{
		ID:         userPlatform.ID.Hex(),
		UserID:     userPlatform.UserID.Hex(),
		PlatformID: userPlatform.PlatformID.Hex(),
		Platform:   platformData,
		AliasName:  userPlatform.AliasName,
		Balance:    utils.Decimal128ToFloat64(userPlatform.Balance),
		IsActive:   userPlatform.IsActive,
		CreatedAt:  userPlatform.CreatedAt,
		UpdatedAt:  userPlatform.UpdatedAt,
		DeletedAt:  userPlatform.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(ctx *gin.Context, userPlatforms []*UserPlatform) []*dto.UserPlatformResponse {
	responses := make([]*dto.UserPlatformResponse, len(userPlatforms))
	for i, userPlatform := range userPlatforms {
		responses[i] = c.mapToResponse(ctx, userPlatform)
	}
	return responses
}

func (c *Controller) mapToDropdownResponse(ctx *gin.Context, userPlatform *UserPlatform) *dto.UserPlatformDropdownResponse {
	var platformData *dto.PlatformData
	platform, err := c.platformRepo.GetPlatformByID(ctx, userPlatform.PlatformID)
	if err == nil && platform != nil {
		platformData = &dto.PlatformData{
			ID:       platform.ID.Hex(),
			Name:     platform.Name,
			Type:     platform.Type,
			IsActive: platform.IsActive,
		}
	}

	return &dto.UserPlatformDropdownResponse{
		ID:        userPlatform.ID.Hex(),
		Platform:  platformData,
		AliasName: userPlatform.AliasName,
		Balance:   utils.Decimal128ToFloat64(userPlatform.Balance),
		IsActive:  userPlatform.IsActive,
	}
}

func (c *Controller) mapToDropdownResponseList(ctx *gin.Context, userPlatforms []*UserPlatform) []*dto.UserPlatformDropdownResponse {
	responses := make([]*dto.UserPlatformDropdownResponse, len(userPlatforms))
	for i, userPlatform := range userPlatforms {
		responses[i] = c.mapToDropdownResponse(ctx, userPlatform)
	}
	return responses
}
