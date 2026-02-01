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

func (c *Controller) GetUserCategories(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userIDStr := userID.(string)
	categories, err := c.service.GetUserCategories(ctx, userIDStr)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	categoryResps := make([]*dto.UserCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResps[i] = c.mapToResponse(category)
	}

	resp := utils.NewSuccessResponse("User categories retrieved successfully", categoryResps)
	ctx.JSON(http.StatusOK, resp)
}

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
