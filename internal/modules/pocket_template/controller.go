package pocket_template

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket_template/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreatePocketTemplate godoc
// @Summary Create a new pocket template
// @Description Create a new pocket template (admin only)
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param request body dto.CreatePocketTemplateRequest true "Pocket template details"
// @Success 201 {object} map[string]interface{} "Pocket template created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Security BearerAuth
// @Router /v1/pocket-templates [post]
func (c *Controller) CreatePocketTemplate(ctx *gin.Context) {
	var req dto.CreatePocketTemplateRequest
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

	template, err := c.service.CreatePocketTemplate(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewSuccessResponse("Pocket template created successfully", templateResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetPocketTemplate godoc
// @Summary Get pocket template by ID
// @Description Get a specific pocket template by ID
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param id path string true "Pocket Template ID"
// @Success 200 {object} map[string]interface{} "Pocket template retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid pocket template ID"
// @Failure 404 {object} map[string]interface{} "Pocket template not found"
// @Security BearerAuth
// @Router /v1/pocket-templates/{id} [get]
func (c *Controller) GetPocketTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	template, err := c.service.GetPocketTemplateByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewSuccessResponse("Pocket template retrieved successfully", templateResp)
	ctx.JSON(http.StatusOK, resp)
}

// UpdatePocketTemplate godoc
// @Summary Update pocket template
// @Description Update a pocket template (admin only)
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param id path string true "Pocket Template ID"
// @Param request body dto.UpdatePocketTemplateRequest true "Pocket template update details"
// @Success 200 {object} map[string]interface{} "Pocket template updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Pocket template not found"
// @Security BearerAuth
// @Router /v1/pocket-templates/{id} [put]
func (c *Controller) UpdatePocketTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto.UpdatePocketTemplateRequest
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

	template, err := c.service.UpdatePocketTemplate(ctx, id, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResp := c.mapToResponse(template)
	resp := utils.NewSuccessResponse("Pocket template updated successfully", templateResp)
	ctx.JSON(http.StatusOK, resp)
}

// DeletePocketTemplate godoc
// @Summary Delete pocket template
// @Description Delete a pocket template (admin only)
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param id path string true "Pocket Template ID"
// @Success 200 {object} map[string]interface{} "Pocket template deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid pocket template ID"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Pocket template not found"
// @Security BearerAuth
// @Router /v1/pocket-templates/{id} [delete]
func (c *Controller) DeletePocketTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeletePocketTemplate(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Pocket template deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// ListPocketTemplates godoc
// @Summary List all pocket templates
// @Description Get a list of all pocket templates with pagination
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Pocket templates retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/pocket-templates [get]
func (c *Controller) ListPocketTemplates(ctx *gin.Context) {
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

	templates, err := c.service.ListPocketTemplates(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templatesResp := c.mapToResponseList(templates)
	resp := utils.NewSuccessResponse("Pocket templates retrieved successfully", templatesResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListActivePocketTemplates godoc
// @Summary List active pocket templates
// @Description Get a list of active pocket templates with pagination
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Active pocket templates retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/pocket-templates/active [get]
func (c *Controller) ListActivePocketTemplates(ctx *gin.Context) {
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

	templates, err := c.service.ListActivePocketTemplates(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templatesResp := c.mapToResponseList(templates)
	resp := utils.NewSuccessResponse("Active pocket templates retrieved successfully", templatesResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListPocketTemplatesByType godoc
// @Summary List pocket templates by type
// @Description Get a list of pocket templates filtered by type with pagination
// @Tags Pocket Templates
// @Accept json
// @Produce json
// @Param type path string true "Pocket Template Type" enums(main,saving,allocation)
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Pocket templates retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /v1/pocket-templates/type/{type} [get]
func (c *Controller) ListPocketTemplatesByType(ctx *gin.Context) {
	templateType := ctx.Param("type")

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

	templates, err := c.service.ListPocketTemplatesByType(ctx, templateType, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templatesResp := c.mapToResponseList(templates)
	resp := utils.NewSuccessResponse("Pocket templates retrieved successfully", templatesResp)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(template *PocketTemplate) *dto.PocketTemplateResponse {
	var categoryID *string
	if template.CategoryID != nil {
		id := template.CategoryID.Hex()
		categoryID = &id
	}
	return &dto.PocketTemplateResponse{
		ID:              template.ID.Hex(),
		Name:            template.Name,
		Type:            template.Type,
		CategoryID:      categoryID,
		Icon:            template.Icon,
		IconColor:       template.IconColor,
		BackgroundColor: template.BackgroundColor,
		IsDefault:       template.IsDefault,
		IsActive:        template.IsActive,
		Order:           template.Order,
		CreatedAt:       template.CreatedAt,
		UpdatedAt:       template.UpdatedAt,
		DeletedAt:       template.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(templates []*PocketTemplate) []dto.PocketTemplateResponse {
	responses := make([]dto.PocketTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = *c.mapToResponse(template)
	}
	return responses
}
