package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/HasanNugroho/coin-be/internal/modules/user/dto"
	"github.com/HasanNugroho/coin-be/pkg/errors"
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
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	user, err := c.service.GetUserByID(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	user, err := c.service.UpdateUser(ctx, userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (c *Controller) ListUsers(ctx *gin.Context) {
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

	users, err := c.service.ListUsers(ctx, limit, skip)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *Controller) CreateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	profile, err := c.service.CreateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, profile)
}

func (c *Controller) GetFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	profile, err := c.service.GetFinancialProfile(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (c *Controller) UpdateFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	var req dto.CreateFinancialProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	profile, err := c.service.UpdateFinancialProfile(ctx, userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (c *Controller) DeleteFinancialProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errors.NewErrorResponse("unauthorized"))
		return
	}

	err := c.service.DeleteFinancialProfile(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "financial profile deleted"})
}

func (c *Controller) CreateRole(ctx *gin.Context) {
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	role, err := c.service.CreateRole(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

func (c *Controller) GetRole(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := c.service.GetRole(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (c *Controller) ListRoles(ctx *gin.Context) {
	roles, err := c.service.ListRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (c *Controller) AssignRoleToUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var req dto.AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	err := c.service.AssignRoleToUser(ctx, userID, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "role assigned"})
}

func (c *Controller) GetUserRoles(ctx *gin.Context) {
	userID := ctx.Param("id")

	roles, err := c.service.GetUserRoles(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (c *Controller) RemoveRoleFromUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	roleID := ctx.Param("role_id")

	err := c.service.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.NewErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "role removed"})
}
