package user_category

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateUserCategory godoc
// @Summary Create a new user category
// @Description Create a new category for the authenticated user
// @Tags User Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateUserCategoryRequest true "Create User Category Request"
// @Success 201 {object} map[string]interface{} "User category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-categories [post]
func (c *Controller) CreateUserCategory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateUserCategoryRequest
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

	userIDStr := userID.(string)
	category, err := c.service.CreateUserCategory(ctx, userIDStr, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewCreatedResponse("User category created successfully", categoryResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetUserCategories godoc
// @Summary Get user categories
// @Description Get a list of user categories with pagination, search, and filters
// @Tags User Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search by name"
// @Param type query string false "Filter by transaction type (income/expense)"
// @Param sort_by query string false "Sort by field (name, created_at, updated_at)"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Success 200 {object} map[string]interface{} "User categories retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-categories [get]
func (c *Controller) GetUserCategories(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	// Get query parameters
	transactionType := ctx.Query("type")
	search := ctx.Query("search")

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(ctx, 10)

	// Parse sorting parameters
	allowedSortFields := []string{"name", "created_at", "updated_at"}
	sorting := utils.ParseSortParams(ctx, allowedSortFields, "created_at")

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}
	var searchFilter *string
	if search != "" {
		searchFilter = &search
	}

	userIDStr := userID.(string)
	categories, total, err := c.service.GetUserCategories(
		ctx,
		userIDStr,
		typeFilter,
		searchFilter,
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

	categoryResps := make([]*dto.UserCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToResponse(category)
	}

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	// Build paginated response
	respData := utils.BuildPaginatedResponse(categoryResps, meta)
	resp := utils.NewSuccessResponse("User categories retrieved successfully", respData)
	ctx.JSON(http.StatusOK, resp)
}

// GetUserCategoryByID godoc
// @Summary Get user category by ID
// @Description Get a specific user category by its ID
// @Tags User Categories
// @Accept json
// @Produce json
// @Param id path string true "User Category ID"
// @Success 200 {object} map[string]interface{} "User category retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User category not found"
// @Security BearerAuth
// @Router /v1/user-categories/{id} [get]
func (c *Controller) GetUserCategoryByID(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")
	userIDStr := userID.(string)

	category, err := c.service.GetUserCategoryByID(ctx, id, userIDStr)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewSuccessResponse("User category retrieved successfully", categoryResp)
	ctx.JSON(http.StatusOK, resp)
}

// FindAllParent godoc
// @Summary Get all parent user categories
// @Description Get a list of all parent user categories (categories without a parent_id)
// @Tags User Categories
// @Accept json
// @Produce json
// @Param type query string false "Filter by transaction type (income/expense)"
// @Success 200 {object} map[string]interface{} "User categories retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User category not found"
// @Security BearerAuth
// @Router /v1/user-categories/parents [get]
func (c *Controller) FindAllParent(ctx *gin.Context) {
	transactionType := ctx.Query("type")
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}

	userIDStr := userID.(string)
	categories, err := c.service.FindAllParent(ctx, userIDStr, typeFilter)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	categoryResps := make([]*dto.SimpleUserCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToSimpleResponse(category)
	}

	resp := utils.NewSuccessResponse("User category retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

// FindAllDropdown godoc
// @Summary Get all user categories for dropdown
// @Description Get a simplified list of all user categories, typically for use in dropdown selectors
// @Tags User Categories
// @Accept json
// @Produce json
// @Param type query string false "Filter by transaction type (income/expense)"
// @Success 200 {object} map[string]interface{} "User categories retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User category not found"
// @Security BearerAuth
// @Router /v1/user-categories/dropdown [get]
func (c *Controller) FindAllDropdown(ctx *gin.Context) {
	transactionType := ctx.Query("type")
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}

	userIDStr := userID.(string)
	categories, err := c.service.FindAllDropdown(ctx, userIDStr, typeFilter)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	categoryResps := make([]*dto.SimpleUserCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToSimpleResponse(category)
	}

	resp := utils.NewSuccessResponse("User category retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateUserCategory godoc
// @Summary Update a user category
// @Description Update an existing user category by ID
// @Tags User Categories
// @Accept json
// @Produce json
// @Param id path string true "User Category ID"
// @Param request body dto.UpdateUserCategoryRequest true "Update User Category Request"
// @Success 200 {object} map[string]interface{} "User category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/user-categories/{id} [put]
func (c *Controller) UpdateUserCategory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	var req dto.UpdateUserCategoryRequest
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

	userIDStr := userID.(string)
	category, err := c.service.UpdateUserCategory(ctx, id, userIDStr, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResp := c.mapToResponse(category)
	resp := utils.NewSuccessResponse("User category updated successfully", categoryResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteUserCategory godoc
// @Summary Delete a user category
// @Description Delete a user category by ID (soft delete)
// @Tags User Categories
// @Accept json
// @Produce json
// @Param id path string true "User Category ID"
// @Success 200 {object} map[string]interface{} "User category deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User category not found"
// @Security BearerAuth
// @Router /v1/user-categories/{id} [delete]
func (c *Controller) DeleteUserCategory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")
	userIDStr := userID.(string)

	err := c.service.DeleteUserCategory(ctx, id, userIDStr)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("User category deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(category *UserCategory) *dto.UserCategoryResponse {
	var parentIDStr *string
	if category.ParentID != nil {
		parentID := category.ParentID.Hex()
		parentIDStr = &parentID
	}

	var templateIDStr *string
	if category.TemplateID != nil {
		templateID := category.TemplateID.Hex()
		templateIDStr = &templateID
	}

	return &dto.UserCategoryResponse{
		ID:              category.ID.Hex(),
		UserID:          category.UserID.Hex(),
		TemplateID:      templateIDStr,
		Name:            category.Name,
		TransactionType: (*string)(category.TransactionType),
		IsDefault:       category.IsDefault,
		Color:           category.Color,
		Icon:            category.Icon,
		Description:     category.Description,
		ParentID:        parentIDStr,
		CreatedAt:       category.CreatedAt,
		UpdatedAt:       category.UpdatedAt,
		DeletedAt:       category.DeletedAt,
	}
}

func (c *Controller) mapToSimpleResponse(category *UserCategory) *dto.SimpleUserCategoryResponse {
	var parentID string
	if category.ParentID != nil {
		parentID = category.ParentID.Hex()
	}

	return &dto.SimpleUserCategoryResponse{
		ID:       category.ID.Hex(),
		UserID:   category.UserID.Hex(),
		ParentID: parentID,
		Name:     category.Name,
		Color:    category.Color,
		Icon:     category.Icon,
	}
}
