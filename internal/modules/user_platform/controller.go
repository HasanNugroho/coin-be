package user_platform

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	userPlatformResp := c.mapToSingleResponse(ctx, userPlatform)
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

	userPlatformResp := c.mapToSingleResponse(ctx, userPlatform)
	resp := utils.NewSuccessResponse("User platform retrieved successfully", userPlatformResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListUserPlatforms godoc
// @Summary List all user platforms
// @Description Get all user platforms for the authenticated user with filtering and pagination
// @Tags User Platforms
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param sort_by query string false "Sort by (alias_name, balance, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param search query string false "Search by alias name"
// @Param is_active query bool false "Filter by active status"
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

	search := ctx.Query("search")
	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	isActiveStr := ctx.Query("is_active")
	var isActivePtr *bool
	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		isActivePtr = &isActive
	}

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(ctx, 10)

	// Parse sorting parameters
	allowedFields := []string{"alias_name", "balance", "created_at"}
	sorting := utils.ParseSortParams(ctx, allowedFields, "created_at")

	userPlatforms, total, err := c.service.ListUserPlatforms(
		ctx,
		userID.(string),
		searchPtr,
		isActivePtr,
		pagination.Page,
		pagination.PageSize,
		sorting.SortBy,
		sorting.SortOrder,
	)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	userPlatformResps := c.mapToResponseList(ctx, userPlatforms)

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	// Build paginated response
	respData := utils.BuildPaginatedResponse(userPlatformResps, meta)
	resp := utils.NewSuccessResponse("User platforms retrieved successfully", respData)
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
// @Router /v1/user-platforms/dropdown [get]
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

	userPlatformResp := c.mapToSingleResponse(ctx, userPlatform)
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

func (c *Controller) buildPlatformMapFromIDs(ctx *gin.Context, platformIDs []primitive.ObjectID) map[primitive.ObjectID]*dto.PlatformData {
	platformMap := make(map[primitive.ObjectID]*dto.PlatformData)
	platforms, err := c.platformRepo.GetPlatformsByIDs(ctx, platformIDs)
	if err != nil {
		return platformMap
	}

	for _, p := range platforms {
		platformMap[p.ID] = &dto.PlatformData{
			ID:       p.ID.Hex(),
			Name:     p.Name,
			Type:     p.Type,
			IsActive: p.IsActive,
		}
	}
	return platformMap
}

func (c *Controller) buildPlatformMap(ctx *gin.Context, userPlatforms []*UserPlatform) map[primitive.ObjectID]*dto.PlatformData {
	platformIDs := make([]primitive.ObjectID, 0, len(userPlatforms))
	for _, up := range userPlatforms {
		platformIDs = append(platformIDs, up.PlatformID)
	}
	return c.buildPlatformMapFromIDs(ctx, platformIDs)
}

func (c *Controller) mapToResponse(ctx *gin.Context, userPlatform *UserPlatform, platformMap map[primitive.ObjectID]*dto.PlatformData) *dto.UserPlatformResponse {
	return &dto.UserPlatformResponse{
		ID:         userPlatform.ID.Hex(),
		UserID:     userPlatform.UserID.Hex(),
		PlatformID: userPlatform.PlatformID.Hex(),
		Platform:   platformMap[userPlatform.PlatformID],
		AliasName:  userPlatform.AliasName,
		Balance:    utils.Decimal128ToFloat64(userPlatform.Balance),
		IsActive:   userPlatform.IsActive,
		CreatedAt:  userPlatform.CreatedAt,
		UpdatedAt:  userPlatform.UpdatedAt,
		DeletedAt:  userPlatform.DeletedAt,
	}
}

func (c *Controller) mapToSingleResponse(ctx *gin.Context, userPlatform *UserPlatform) *dto.UserPlatformResponse {
	platformMap := c.buildPlatformMapFromIDs(ctx, []primitive.ObjectID{userPlatform.PlatformID})
	return c.mapToResponse(ctx, userPlatform, platformMap)
}

func (c *Controller) mapToResponseList(ctx *gin.Context, userPlatforms []*UserPlatform) []*dto.UserPlatformResponse {
	platformMap := c.buildPlatformMap(ctx, userPlatforms)
	responses := make([]*dto.UserPlatformResponse, len(userPlatforms))
	for i, up := range userPlatforms {
		responses[i] = c.mapToResponse(ctx, up, platformMap)
	}
	return responses
}

func (c *Controller) mapToDropdownResponseList(ctx *gin.Context, userPlatforms []*UserPlatform) []*dto.UserPlatformDropdownResponse {
	platformMap := c.buildPlatformMap(ctx, userPlatforms)
	responses := make([]*dto.UserPlatformDropdownResponse, len(userPlatforms))
	for i, up := range userPlatforms {
		responses[i] = &dto.UserPlatformDropdownResponse{
			ID:        up.ID.Hex(),
			Platform:  platformMap[up.PlatformID],
			AliasName: up.AliasName,
			Balance:   utils.Decimal128ToFloat64(up.Balance),
			IsActive:  up.IsActive,
		}
	}
	return responses
}
