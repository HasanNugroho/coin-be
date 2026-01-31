package platform

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/platform/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreatePlatform godoc
// @Summary Create a new platform
// @Description Create a new platform (admin only)
// @Tags Platforms
// @Accept json
// @Produce json
// @Param request body dto.CreatePlatformRequest true "Platform details"
// @Success 201 {object} map[string]interface{} "Platform created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /v1/platforms [post]
func (c *Controller) CreatePlatform(ctx *gin.Context) {
	var req dto.CreatePlatformRequest
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

	platform, err := c.service.CreatePlatform(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	platformResp := c.mapToResponse(platform)
	resp := utils.NewSuccessResponse("Platform created successfully", platformResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetPlatform godoc
// @Summary Get platform by ID
// @Description Get a specific platform by ID
// @Tags Platforms
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object} map[string]interface{} "Platform retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid platform ID"
// @Failure 404 {object} map[string]interface{} "Platform not found"
// @Security BearerAuth
// @Router /v1/platforms/{id} [get]
func (c *Controller) GetPlatform(ctx *gin.Context) {
	id := ctx.Param("id")

	platform, err := c.service.GetPlatformByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	platformResp := c.mapToResponse(platform)
	resp := utils.NewSuccessResponse("Platform retrieved successfully", platformResp)
	ctx.JSON(http.StatusOK, resp)
}

// UpdatePlatform godoc
// @Summary Update platform
// @Description Update a platform (admin only)
// @Tags Platforms
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Param request body dto.UpdatePlatformRequest true "Platform update details"
// @Success 200 {object} map[string]interface{} "Platform updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Platform not found"
// @Security BearerAuth
// @Router /v1/platforms/{id} [put]
func (c *Controller) UpdatePlatform(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdatePlatformRequest
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

	platform, err := c.service.UpdatePlatform(ctx, id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	platformResp := c.mapToResponse(platform)
	resp := utils.NewSuccessResponse("Platform updated successfully", platformResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeletePlatform godoc
// @Summary Delete platform
// @Description Soft delete a platform (admin only)
// @Tags Platforms
// @Accept json
// @Produce json
// @Param id path string true "Platform ID"
// @Success 200 {object} map[string]interface{} "Platform deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Platform not found"
// @Security BearerAuth
// @Router /v1/platforms/{id} [delete]
func (c *Controller) DeletePlatform(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeletePlatform(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Platform deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// ListPlatforms godoc
// @Summary List all platforms
// @Description Get a list of all platforms with pagination
// @Tags Platforms
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Platforms retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/platforms [get]
func (c *Controller) ListPlatforms(ctx *gin.Context) {
	limit := int64(10)
	skip := int64(0)

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}

	if s := ctx.Query("skip"); s != "" {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			skip = parsed
		}
	}

	platforms, err := c.service.ListPlatforms(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	platformsResp := c.mapToResponseList(platforms)
	resp := utils.NewSuccessResponse("Platforms retrieved successfully", platformsResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListActivePlatforms godoc
// @Summary List active platforms
// @Description Get a list of active platforms with pagination
// @Tags Platforms
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Active platforms retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/platforms/active [get]
func (c *Controller) ListActivePlatforms(ctx *gin.Context) {
	limit := int64(10)
	skip := int64(0)

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}

	if s := ctx.Query("skip"); s != "" {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			skip = parsed
		}
	}

	platforms, err := c.service.ListActivePlatforms(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	platformsResp := c.mapToResponseList(platforms)
	resp := utils.NewSuccessResponse("Active platforms retrieved successfully", platformsResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListPlatformsByType godoc
// @Summary List platforms by type
// @Description Get a list of platforms filtered by type with pagination
// @Tags Platforms
// @Accept json
// @Produce json
// @Param type query string true "Platform type (BANK, E_WALLET, CASH, ATM)"
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Platforms retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/platforms/type/{type} [get]
func (c *Controller) ListPlatformsByType(ctx *gin.Context) {
	platformType := ctx.Param("type")

	limit := int64(10)
	skip := int64(0)

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}

	if s := ctx.Query("skip"); s != "" {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			skip = parsed
		}
	}

	platforms, err := c.service.ListPlatformsByType(ctx, platformType, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	platformsResp := c.mapToResponseList(platforms)
	resp := utils.NewSuccessResponse("Platforms retrieved successfully", platformsResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(platform *Platform) *dto.PlatformResponse {
	return &dto.PlatformResponse{
		ID:        platform.ID.Hex(),
		Name:      platform.Name,
		Type:      platform.Type,
		IsActive:  platform.IsActive,
		CreatedAt: platform.CreatedAt,
		UpdatedAt: platform.UpdatedAt,
		DeletedAt: platform.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(platforms []*Platform) []*dto.PlatformResponse {
	responses := make([]*dto.PlatformResponse, len(platforms))
	for i, platform := range platforms {
		responses[i] = c.mapToResponse(platform)
	}
	return responses
}
