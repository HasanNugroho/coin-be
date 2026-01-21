package allocation

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateAllocation godoc
// @Summary Create a new allocation
// @Description Create a new allocation with priority and percentage
// @Tags Allocations
// @Accept json
// @Produce json
// @Param request body dto.CreateAllocationRequest true "Allocation details"
// @Success 201 {object} map[string]interface{} "Allocation created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations [post]
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

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocation, err := c.service.CreateAllocation(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation created successfully", allocation)
	ctx.JSON(http.StatusCreated, resp)
}

// GetAllocations godoc
// @Summary Get all allocations
// @Description Get all allocations for the authenticated user
// @Tags Allocations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Allocations retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations [get]
func (c *Controller) GetAllocations(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocations, err := c.service.GetAllocations(ctx, userOID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocations retrieved successfully", allocations)
	ctx.JSON(http.StatusOK, resp)
}

// GetAllocationByID godoc
// @Summary Get allocation by ID
// @Description Get a specific allocation by its ID
// @Tags Allocations
// @Accept json
// @Produce json
// @Param id path string true "Allocation ID"
// @Success 200 {object} map[string]interface{} "Allocation retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid allocation ID"
// @Failure 404 {object} map[string]interface{} "Allocation not found"
// @Security BearerAuth
// @Router /allocations/{id} [get]
func (c *Controller) GetAllocationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	allocationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid allocation ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocation, err := c.service.GetAllocationByID(ctx, allocationID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation retrieved successfully", allocation)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateAllocation godoc
// @Summary Update allocation
// @Description Update an existing allocation
// @Tags Allocations
// @Accept json
// @Produce json
// @Param id path string true "Allocation ID"
// @Param request body dto.UpdateAllocationRequest true "Updated allocation details"
// @Success 200 {object} map[string]interface{} "Allocation updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations/{id} [put]
func (c *Controller) UpdateAllocation(ctx *gin.Context) {
	id := ctx.Param("id")
	allocationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid allocation ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req dto.UpdateAllocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	allocation, err := c.service.UpdateAllocation(ctx, allocationID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation updated successfully", allocation)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteAllocation godoc
// @Summary Delete allocation
// @Description Delete an allocation by ID
// @Tags Allocations
// @Accept json
// @Produce json
// @Param id path string true "Allocation ID"
// @Success 200 {object} map[string]interface{} "Allocation deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid allocation ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations/{id} [delete]
func (c *Controller) DeleteAllocation(ctx *gin.Context) {
	id := ctx.Param("id")
	allocationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid allocation ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.service.DeleteAllocation(ctx, allocationID); err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// GetAllocationLogs godoc
// @Summary Get allocation distribution logs
// @Description Get all allocation distribution logs for the authenticated user
// @Tags Allocations
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param skip query int false "Skip" default(0)
// @Success 200 {object} map[string]interface{} "Allocation logs retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations/logs [get]
func (c *Controller) GetAllocationLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "50"), 10, 64)
	skip, _ := strconv.ParseInt(ctx.DefaultQuery("skip", "0"), 10, 64)

	logs, err := c.service.GetAllocationLogs(ctx, userOID, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation logs retrieved successfully", logs)
	ctx.JSON(http.StatusOK, resp)
}

// GetAllocationLogsByAllocationID godoc
// @Summary Get allocation logs by allocation ID
// @Description Get distribution logs for a specific allocation
// @Tags Allocations
// @Accept json
// @Produce json
// @Param id path string true "Allocation ID"
// @Param limit query int false "Limit" default(50)
// @Param skip query int false "Skip" default(0)
// @Success 200 {object} map[string]interface{} "Allocation logs retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid allocation ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /allocations/{id}/logs [get]
func (c *Controller) GetAllocationLogsByAllocationID(ctx *gin.Context) {
	id := ctx.Param("id")
	allocationID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid allocation ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "50"), 10, 64)
	skip, _ := strconv.ParseInt(ctx.DefaultQuery("skip", "0"), 10, 64)

	logs, err := c.service.GetAllocationLogsByAllocationID(ctx, allocationID, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation logs retrieved successfully", logs)
	ctx.JSON(http.StatusOK, resp)
}
