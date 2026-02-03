package allocation

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateAllocation godoc
// @Summary Create a new allocation
// @Description Create a new salary allocation rule
// @Tags Allocations
// @Accept json
// @Produce json
// @Param request body dto.CreateAllocationRequest true "Allocation details"
// @Success 201 {object} map[string]interface{} "Allocation created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/allocations [post]
func (c *Controller) CreateAllocation(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateAllocationRequest
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

	allocation, err := c.service.CreateAllocation(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocationResp := c.mapToResponse(allocation)
	resp := utils.NewSuccessResponse("Allocation created successfully", allocationResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetAllocation godoc
// @Summary Get allocation by ID
// @Description Get a specific allocation by ID
// @Tags Allocations
// @Produce json
// @Param id path string true "Allocation ID"
// @Success 200 {object} map[string]interface{} "Allocation retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/allocations/{id} [get]
func (c *Controller) GetAllocation(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	allocation, err := c.service.GetAllocationByID(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocationResp := c.mapToResponse(allocation)
	resp := utils.NewSuccessResponse("Allocation retrieved successfully", allocationResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListAllocations godoc
// @Summary List all allocations
// @Description Get all allocations for the authenticated user
// @Tags Allocations
// @Produce json
// @Success 200 {object} map[string]interface{} "Allocations retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/allocations [get]
func (c *Controller) ListAllocations(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	allocations, err := c.service.ListAllocations(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocationResps := c.mapToResponseList(allocations)
	resp := utils.NewSuccessResponse("Allocations retrieved successfully", allocationResps)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateAllocation godoc
// @Summary Update allocation
// @Description Update an allocation
// @Tags Allocations
// @Accept json
// @Produce json
// @Param id path string true "Allocation ID"
// @Param request body dto.UpdateAllocationRequest true "Update details"
// @Success 200 {object} map[string]interface{} "Allocation updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/allocations/{id} [put]
func (c *Controller) UpdateAllocation(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	var req dto.UpdateAllocationRequest
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

	allocation, err := c.service.UpdateAllocation(ctx, userID.(string), id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocationResp := c.mapToResponse(allocation)
	resp := utils.NewSuccessResponse("Allocation updated successfully", allocationResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteAllocation godoc
// @Summary Delete allocation
// @Description Delete an allocation
// @Tags Allocations
// @Produce json
// @Param id path string true "Allocation ID"
// @Success 200 {object} map[string]interface{} "Allocation deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Security BearerAuth
// @Router /v1/allocations/{id} [delete]
func (c *Controller) DeleteAllocation(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	err := c.service.DeleteAllocation(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(allocation *Allocation) *dto.AllocationResponse {
	var pocketID *string
	if allocation.PocketID != nil {
		id := allocation.PocketID.Hex()
		pocketID = &id
	}

	var userPlatformID *string
	if allocation.UserPlatformID != nil {
		id := allocation.UserPlatformID.Hex()
		userPlatformID = &id
	}

	return &dto.AllocationResponse{
		ID:             allocation.ID.Hex(),
		UserID:         allocation.UserID.Hex(),
		PocketID:       pocketID,
		UserPlatformID: userPlatformID,
		Priority:       allocation.Priority,
		AllocationType: allocation.AllocationType,
		Nominal:        allocation.Nominal,
		IsActive:       allocation.IsActive,
		ExecuteDay:     allocation.ExecuteDay,
		CreatedAt:      allocation.CreatedAt,
		UpdatedAt:      allocation.UpdatedAt,
		DeletedAt:      allocation.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(allocations []*Allocation) []*dto.AllocationResponse {
	responses := make([]*dto.AllocationResponse, len(allocations))
	for i, allocation := range allocations {
		responses[i] = c.mapToResponse(allocation)
	}
	return responses
}
