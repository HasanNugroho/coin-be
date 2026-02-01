package category_template

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/category_template/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (c *Controller) CreateCategoryTemplate(ctx *gin.Context) {
	var req dto.CreateCategoryTemplateRequest
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

	template, err := c.service.CreateCategoryTemplate(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewCreatedResponse("Category template created successfully", templateResp)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) FindAll(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	// Get query parameters
	transactionType := ctx.Query("type")
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")

	// Parse pagination parameters
	var pageNum int64 = 1
	var pageSizeNum int64 = 10

	if p, err := strconv.ParseInt(page, 10, 64); err == nil && p > 0 {
		pageNum = p
	}
	if ps, err := strconv.ParseInt(pageSize, 10, 64); err == nil && ps > 0 {
		pageSizeNum = ps
	}

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}

	// Fetch templates with filter and pagination
	templates, total, err := c.service.FindAllWithFilter(ctx, typeFilter, pageNum, pageSizeNum)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResps := make([]*dto.CategoryTemplateResponse, len(templates))
	for i, template := range templates {
		templateResps[i] = c.mapToResponse(template)
	}

	// Calculate pagination metadata
	totalPages := (total + pageSizeNum - 1) / pageSizeNum
	meta := gin.H{
		"total":       total,
		"page":        pageNum,
		"page_size":   pageSizeNum,
		"total_pages": totalPages,
	}

	resp := utils.NewSuccessResponse("Category templates retrieved successfully", gin.H{
		"data": templateResps,
		"meta": meta,
	})
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) FindAllParent(ctx *gin.Context) {
	_, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "user id not found in context")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	transactionType := ctx.Query("type")
	templates, err := c.service.FindAllParent(ctx, &transactionType)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category templates retrieved successfully", templates)
	ctx.JSON(http.StatusOK, resp)

}

func (c *Controller) GetCategoryTemplateByID(ctx *gin.Context) {
	id := ctx.Param("id")

	template, err := c.service.GetCategoryTemplateByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewSuccessResponse("Category template retrieved successfully", templateResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UpdateCategoryTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdateCategoryTemplateRequest
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

	template, err := c.service.UpdateCategoryTemplate(ctx, id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewSuccessResponse("Category template updated successfully", templateResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) DeleteCategoryTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteCategoryTemplate(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Category template deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(template *CategoryTemplate) *dto.CategoryTemplateResponse {
	var parentIDStr *string
	if template.ParentID != nil {
		parentID := template.ParentID.Hex()
		parentIDStr = &parentID
	}

	var userIDStr *string
	if template.UserID != nil {
		userID := template.UserID.Hex()
		userIDStr = &userID
	}

	return &dto.CategoryTemplateResponse{
		ID:              template.ID.Hex(),
		Name:            template.Name,
		TransactionType: (*string)(template.TransactionType),
		IsDefault:       template.IsDefault,
		Color:           template.Color,
		Icon:            template.Icon,
		Description:     template.Description,
		ParentID:        parentIDStr,
		UserID:          userIDStr,
		CreatedAt:       template.CreatedAt,
		UpdatedAt:       template.UpdatedAt,
		DeletedAt:       template.DeletedAt,
	}
}
