package target

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/target/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateTarget godoc
// @Summary Create a new saving target
// @Description Create a new saving target linked to an allocation
// @Tags Saving Targets
// @Accept json
// @Produce json
// @Param request body dto.CreateTargetRequest true "Target details"
// @Success 201 {object} map[string]interface{} "Target created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /targets [post]
func (c *Controller) CreateTarget(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateTargetRequest
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

	target, err := c.service.CreateTarget(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Target created successfully", target)
	ctx.JSON(http.StatusCreated, resp)
}

// GetTargets godoc
// @Summary Get all saving targets
// @Description Get all saving targets for the authenticated user
// @Tags Saving Targets
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Targets retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /targets [get]
func (c *Controller) GetTargets(ctx *gin.Context) {
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

	targets, err := c.service.GetTargets(ctx, userOID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Targets retrieved successfully", targets)
	ctx.JSON(http.StatusOK, resp)
}

// GetTargetByID godoc
// @Summary Get target by ID
// @Description Get a specific saving target by its ID
// @Tags Saving Targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Success 200 {object} map[string]interface{} "Target retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid target ID"
// @Failure 404 {object} map[string]interface{} "Target not found"
// @Security BearerAuth
// @Router /targets/{id} [get]
func (c *Controller) GetTargetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	targetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid target ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	target, err := c.service.GetTargetByID(ctx, targetID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Target retrieved successfully", target)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateTarget godoc
// @Summary Update saving target
// @Description Update an existing saving target
// @Tags Saving Targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Param request body dto.UpdateTargetRequest true "Updated target details"
// @Success 200 {object} map[string]interface{} "Target updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /targets/{id} [put]
func (c *Controller) UpdateTarget(ctx *gin.Context) {
	id := ctx.Param("id")
	targetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid target ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req dto.UpdateTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	target, err := c.service.UpdateTarget(ctx, targetID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Target updated successfully", target)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteTarget godoc
// @Summary Delete saving target
// @Description Delete a saving target by ID
// @Tags Saving Targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Success 200 {object} map[string]interface{} "Target deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid target ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /targets/{id} [delete]
func (c *Controller) DeleteTarget(ctx *gin.Context) {
	id := ctx.Param("id")
	targetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid target ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.service.DeleteTarget(ctx, targetID); err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Target deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
