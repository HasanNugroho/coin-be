package category

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/category/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new income or expense category
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories [post]
func (c *Controller) CreateCategory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateCategoryRequest
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

	category, err := c.service.CreateCategory(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category created successfully", category)
	ctx.JSON(http.StatusCreated, resp)
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories for the authenticated user. Filter by type using query parameter
// @Tags Categories
// @Accept json
// @Produce json
// @Param type query string false "Category type (income or expense)"
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories [get]
func (c *Controller) GetCategories(ctx *gin.Context) {
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

	categoryType := ctx.Query("type")
	var categories []*dto.CategoryResponse

	if categoryType != "" {
		categories, err = c.service.GetCategoriesByType(ctx, userOID, categoryType)
	} else {
		categories, err = c.service.GetCategories(ctx, userOID)
	}

	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Categories retrieved successfully", categories)
	ctx.JSON(http.StatusOK, resp)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /categories/{id} [get]
func (c *Controller) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid category ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	category, err := c.service.GetCategoryByID(ctx, categoryID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category retrieved successfully", category)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Updated category details"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories/{id} [put]
func (c *Controller) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid category ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	category, err := c.service.UpdateCategory(ctx, categoryID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category updated successfully", category)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (c *Controller) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid category ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.service.DeleteCategory(ctx, categoryID); err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
