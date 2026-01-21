package user

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/user/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (c *Controller) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	user, err := c.service.GetUserByID(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Profile retrieved successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := c.service.UpdateUser(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Profile updated successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetUserByID(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("User retrieved successfully", user)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteUser(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) ListUsers(ctx *gin.Context) {
	page := int64(1)
	limit := int64(10)

	if p := ctx.Query("page"); p != "" {
		if parsed, err := strconv.ParseInt(p, 10, 64); err == nil {
			page = parsed
		}
	}

	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsed
		}
	}

	skip := (page - 1) * limit
	users, err := c.service.ListUsers(ctx, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	total := int64(len(users))
	pagination := utils.CalculatePagination(page, limit, total)
	resp := utils.NewSuccessResponseWithPagination("Users retrieved successfully", users, pagination)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) CreateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	profile, err := c.service.CreateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewCreatedResponse("Financial profile created successfully", profile)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) GetFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	profile, err := c.service.GetFinancialProfile(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile retrieved successfully", profile)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UpdateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	profile, err := c.service.UpdateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile updated successfully", profile)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) DeleteFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	err := c.service.DeleteFinancialProfile(ctx, userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Financial profile deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) CreateRole(ctx *gin.Context) {
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	role, err := c.service.CreateRole(ctx, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewCreatedResponse("Role created successfully", role)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *Controller) GetRole(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := c.service.GetRole(ctx, id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Role retrieved successfully", role)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) ListRoles(ctx *gin.Context) {
	roles, err := c.service.ListRoles(ctx)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Roles retrieved successfully", roles)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) AssignRoleToUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var req dto.AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	err := c.service.AssignRoleToUser(ctx, userID, req.RoleID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Role assigned successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetUserRoles(ctx *gin.Context) {
	userID := ctx.Param("id")

	roles, err := c.service.GetUserRoles(ctx, userID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("User roles retrieved successfully", roles)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) RemoveRoleFromUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	roleID := ctx.Param("role_id")

	err := c.service.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Role removed successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
