package pocket

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreatePocket godoc
// @Summary Create a new pocket
// @Description Create a new pocket for the authenticated user
// @Tags Pockets
// @Accept json
// @Produce json
// @Param request body dto.CreatePocketRequest true "Pocket details"
// @Success 201 {object} map[string]interface{} "Pocket created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/pockets [post]
func (c *Controller) CreatePocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreatePocketRequest
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

	pocket, err := c.service.CreatePocket(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketResp := c.mapToResponse(pocket)
	resp := utils.NewSuccessResponse("Pocket created successfully", pocketResp)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) LockPocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	c.service.ToggleLockPocket(ctx, userID.(string), id, true)

	resp := utils.NewSuccessResponse("Pocket locked successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UnlockPocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	c.service.ToggleLockPocket(ctx, userID.(string), id, false)

	resp := utils.NewSuccessResponse("Pocket unlocked successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// GetPocket godoc
// @Summary Get pocket by ID
// @Description Get a specific pocket by ID
// @Tags Pockets
// @Accept json
// @Produce json
// @Param id path string true "Pocket ID"
// @Success 200 {object} map[string]interface{} "Pocket retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Pocket not found"
// @Security BearerAuth
// @Router /v1/pockets/{id} [get]
func (c *Controller) GetPocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	pocket, err := c.service.GetPocketByID(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	pocketResp := c.mapToResponse(pocket)
	resp := utils.NewSuccessResponse("Pocket retrieved successfully", pocketResp)
	ctx.JSON(http.StatusOK, resp)
}

// UpdatePocket godoc
// @Summary Update pocket
// @Description Update a pocket
// @Tags Pockets
// @Accept json
// @Produce json
// @Param id path string true "Pocket ID"
// @Param request body dto.UpdatePocketRequest true "Pocket update details"
// @Success 200 {object} map[string]interface{} "Pocket updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Pocket not found"
// @Security BearerAuth
// @Router /v1/pockets/{id} [put]
func (c *Controller) UpdatePocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	var req dto.UpdatePocketRequest
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

	pocket, err := c.service.UpdatePocket(ctx, userID.(string), id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketResp := c.mapToResponse(pocket)
	resp := utils.NewSuccessResponse("Pocket updated successfully", pocketResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeletePocket godoc
// @Summary Delete pocket
// @Description Soft delete a pocket
// @Tags Pockets
// @Accept json
// @Produce json
// @Param id path string true "Pocket ID"
// @Success 200 {object} map[string]interface{} "Pocket deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Pocket not found"
// @Security BearerAuth
// @Router /v1/pockets/{id} [delete]
func (c *Controller) DeletePocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	err := c.service.DeletePocket(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Pocket deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// ListPockets godoc
// @Summary List user pockets
// @Description Get a list of user pockets with pagination
// @Tags Pockets
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Pockets retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/pockets [get]
func (c *Controller) ListPockets(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

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

	pockets, err := c.service.GetUserPockets(ctx, userID.(string), limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketsResp := c.mapToResponseList(pockets)
	resp := utils.NewSuccessResponse("Pockets retrieved successfully", pocketsResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListActivePockets godoc
// @Summary List active user pockets
// @Description Get a list of active user pockets with pagination
// @Tags Pockets
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Active pockets retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/pockets/active [get]
func (c *Controller) ListActivePockets(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

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

	pockets, err := c.service.GetActiveUserPockets(ctx, userID.(string), limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketsResp := c.mapToResponseList(pockets)
	resp := utils.NewSuccessResponse("Active pockets retrieved successfully", pocketsResp)
	ctx.JSON(http.StatusOK, resp)
}

// GetMainPocket godoc
// @Summary Get main pocket
// @Description Get the user's main pocket
// @Tags Pockets
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Main pocket retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Main pocket not found"
// @Security BearerAuth
// @Router /v1/pockets/main [get]
func (c *Controller) GetMainPocket(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	pocket, err := c.service.GetMainPocket(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	pocketResp := c.mapToResponse(pocket)
	resp := utils.NewSuccessResponse("Main pocket retrieved successfully", pocketResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListPocketsDropdown godoc
// @Summary List pockets for dropdown
// @Description Get all pockets for dropdown with platform data lookup
// @Tags Pockets
// @Produce json
// @Success 200 {object} map[string]interface{} "Pockets retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/pockets/dropdown [get]
func (c *Controller) ListPocketsDropdown(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	pockets, err := c.service.ListPocketsDropdown(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketResps := c.mapToResponseDropdownList(pockets)
	resp := utils.NewSuccessResponse("Pockets retrieved successfully", pocketResps)
	ctx.JSON(http.StatusOK, resp)
}

// Admin endpoints

// CreateSystemPocket godoc
// @Summary Create system pocket (admin only)
// @Description Create a system pocket for a user (admin only)
// @Tags Pockets Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body dto.CreateSystemPocketRequest true "System pocket details"
// @Success 201 {object} map[string]interface{} "System pocket created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /v1/admin/pockets/{user_id} [post]
func (c *Controller) CreateSystemPocket(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	var req dto.CreateSystemPocketRequest
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

	pocket, err := c.service.CreateSystemPocket(ctx, userID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketResp := c.mapToResponse(pocket)
	resp := utils.NewSuccessResponse("System pocket created successfully", pocketResp)
	ctx.JSON(http.StatusCreated, resp)
}

// ListAllPockets godoc
// @Summary List all pockets (admin only)
// @Description Get a list of all pockets with pagination (admin only)
// @Tags Pockets Admin
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Pockets retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /v1/admin/pockets [get]
func (c *Controller) ListAllPockets(ctx *gin.Context) {
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

	pockets, err := c.service.GetAllPockets(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	pocketsResp := c.mapToResponseList(pockets)
	resp := utils.NewSuccessResponse("Pockets retrieved successfully", pocketsResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(pocket *Pocket) *dto.PocketResponse {
	var categoryID *string
	if pocket.CategoryID != nil {
		id := pocket.CategoryID.Hex()
		categoryID = &id
	}

	var targetBalance *float64
	if pocket.TargetBalance != nil {
		tb := utils.Decimal128ToFloat64(*pocket.TargetBalance)
		targetBalance = &tb
	}

	return &dto.PocketResponse{
		ID:              pocket.ID.Hex(),
		UserID:          pocket.UserID.Hex(),
		Name:            pocket.Name,
		Type:            pocket.Type,
		CategoryID:      categoryID,
		Balance:         utils.Decimal128ToFloat64(pocket.Balance),
		TargetBalance:   targetBalance,
		IsDefault:       pocket.IsDefault,
		IsActive:        pocket.IsActive,
		IsLocked:        pocket.IsLocked,
		Icon:            pocket.Icon,
		IconColor:       pocket.IconColor,
		BackgroundColor: pocket.BackgroundColor,
		CreatedAt:       pocket.CreatedAt,
		UpdatedAt:       pocket.UpdatedAt,
		DeletedAt:       pocket.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(pockets []*Pocket) []*dto.PocketResponse {
	responses := make([]*dto.PocketResponse, len(pockets))
	for i, pocket := range pockets {
		responses[i] = c.mapToResponse(pocket)
	}
	return responses
}

func (c *Controller) mapToResponseDropdown(pocket *Pocket) *dto.PocketDropdownResponse {
	return &dto.PocketDropdownResponse{
		ID:              pocket.ID.Hex(),
		Name:            pocket.Name,
		Type:            pocket.Type,
		Balance:         utils.Decimal128ToFloat64(pocket.Balance),
		BackgroundColor: pocket.BackgroundColor,
		IsActive:        pocket.IsActive,
		IsLocked:        pocket.IsLocked,
	}
}

func (c *Controller) mapToResponseDropdownList(pockets []*Pocket) []*dto.PocketDropdownResponse {
	responses := make([]*dto.PocketDropdownResponse, len(pockets))
	for i, pocket := range pockets {
		responses[i] = c.mapToResponseDropdown(pocket)
	}
	return responses
}
