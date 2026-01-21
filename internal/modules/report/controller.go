package report

import (
	"net/http"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// GetDashboardSummary godoc
// @Summary Get dashboard summary
// @Description Get financial dashboard summary with total balance, income, and expenses
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Dashboard summary retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /reports/dashboard [get]
func (c *Controller) GetDashboardSummary(ctx *gin.Context) {
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

	summary, err := c.service.GetDashboardSummary(ctx, userOID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Dashboard summary retrieved successfully", summary)
	ctx.JSON(http.StatusOK, resp)
}

// GetIncomeReport godoc
// @Summary Get income report
// @Description Get detailed income report with breakdown by category and month
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Income report retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /reports/income [get]
func (c *Controller) GetIncomeReport(ctx *gin.Context) {
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

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid start_date format")
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid end_date format")
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	} else {
		endDate = time.Now()
	}

	report, err := c.service.GetIncomeReport(ctx, userOID, startDate, endDate)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Income report retrieved successfully", report)
	ctx.JSON(http.StatusOK, resp)
}

// GetExpenseReport godoc
// @Summary Get expense report
// @Description Get detailed expense report with breakdown by category, allocation, and month
// @Tags Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Expense report retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /reports/expense [get]
func (c *Controller) GetExpenseReport(ctx *gin.Context) {
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

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid start_date format")
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid end_date format")
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
	} else {
		endDate = time.Now()
	}

	report, err := c.service.GetExpenseReport(ctx, userOID, startDate, endDate)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Expense report retrieved successfully", report)
	ctx.JSON(http.StatusOK, resp)
}

// GetAllocationReport godoc
// @Summary Get allocation report
// @Description Get detailed allocation report with balances and distribution history
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Allocation report retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /reports/allocation [get]
func (c *Controller) GetAllocationReport(ctx *gin.Context) {
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

	report, err := c.service.GetAllocationReport(ctx, userOID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Allocation report retrieved successfully", report)
	ctx.JSON(http.StatusOK, resp)
}

// GetTargetProgress godoc
// @Summary Get target progress
// @Description Get progress report for all saving targets
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Target progress retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /reports/target-progress [get]
func (c *Controller) GetTargetProgress(ctx *gin.Context) {
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

	progress, err := c.service.GetTargetProgress(ctx, userOID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Target progress retrieved successfully", progress)
	ctx.JSON(http.StatusOK, resp)
}
