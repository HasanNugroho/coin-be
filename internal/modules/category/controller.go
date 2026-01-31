package category

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/category/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new transaction or pocket category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /v1/categories [post]
func (c *Controller) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	category, err := c.service.CreateCategory(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewSuccessResponse("Category created successfully", categoryResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Get a specific category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /v1/categories/{id} [get]
func (c *Controller) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := c.service.GetCategoryByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewSuccessResponse("Category retrieved successfully", categoryResp)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update a category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /v1/categories/{id} [put]
func (c *Controller) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	category, err := c.service.UpdateCategory(ctx, id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewSuccessResponse("Category updated successfully", categoryResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Soft delete a category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /v1/categories/{id} [delete]
func (c *Controller) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteCategory(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// ListCategories godoc
// @Summary List all categories
// @Description Get a paginated list of all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/categories [get]
func (c *Controller) ListCategories(ctx *gin.Context) {
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

	categories, err := c.service.ListCategories(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResps := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToResponse(category)
	}

	resp := utils.NewSuccessResponse("Categories retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

// ListCategoriesByType godoc
// @Summary List categories by type
// @Description Get categories filtered by type (transaction or pocket)
// @Tags Categories
// @Accept json
// @Produce json
// @Param type query string true "Category type (transaction or pocket)"
// @Param limit query int false "Limit (default: 10)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/categories/type/{type} [get]
func (c *Controller) ListCategoriesByType(ctx *gin.Context) {
	categoryType := ctx.Param("type")
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

	categories, err := c.service.ListCategoriesByType(ctx, categoryType, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResps := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToResponse(category)
	}

	resp := utils.NewSuccessResponse("Categories retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

// ListSubcategories godoc
// @Summary List subcategories
// @Description Get all subcategories of a parent category
// @Tags Categories
// @Accept json
// @Produce json
// @Param parent_id path string true "Parent Category ID"
// @Success 200 {object} map[string]interface{} "Subcategories retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid parent ID"
// @Security BearerAuth
// @Router /v1/categories/{parent_id}/subcategories [get]
func (c *Controller) ListSubcategories(ctx *gin.Context) {
	parentID := ctx.Param("parent_id")

	categories, err := c.service.ListSubcategories(ctx, parentID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResps := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToResponse(category)
	}

	resp := utils.NewSuccessResponse("Subcategories retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(category *Category) *dto.CategoryResponse {
	var parentIDStr *string
	if category.ParentID != nil {
		parentID := category.ParentID.Hex()
		parentIDStr = &parentID
	}

	var userIDStr *string
	if category.UserID != nil {
		userID := category.UserID.Hex()
		userIDStr = &userID
	}

	return &dto.CategoryResponse{
		ID:              category.ID.Hex(),
		Name:            category.Name,
		Type:            string(category.Type),
		TransactionType: (*string)(category.TransactionType),
		IsDefault:       category.IsDefault,
		Color:           category.Color,
		Icon:            category.Icon,
		Description:     category.Description,
		ParentID:        parentIDStr,
		UserID:          userIDStr,
		CreatedAt:       category.CreatedAt,
		UpdatedAt:       category.UpdatedAt,
		DeletedAt:       category.DeletedAt,
	}
}
