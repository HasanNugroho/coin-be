package dashboard

import (
	"net/http"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/daily_summary"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service             *Service
	dailySummaryService *daily_summary.Service
}

func NewController(s *Service, dss *daily_summary.Service) *Controller {
	return &Controller{
		service:             s,
		dailySummaryService: dss,
	}
}

// GetDashboardSummary godoc
// @Summary Get dashboard summary
// @Description Get real-time dashboard summary with total net worth and period income/expense using Hybrid Logic. Default is rolling 30 days. Use filter for 7d (rolling 7 days), 1m (calendar month from 1st), 3m (calendar 3 months from 1st)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param time_range query string false "Time range filter" Enums(7d, 1m, 3m)
// @Success 200 {object} map[string]interface{} "Dashboard summary retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/dashboard/summary [get]
func (c *Controller) GetDashboardSummary(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	timeRange := TimeRange(ctx.DefaultQuery("time_range", ""))
	if timeRange != "" && timeRange != TimeRange7Days && timeRange != TimeRange1Month && timeRange != TimeRange3Month {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid time_range, allowed values: 7d, 1m, 3m")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	summary, err := c.service.GetDashboardSummary(ctx, userID.(string), timeRange)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Dashboard summary retrieved successfully", summary)
	ctx.JSON(http.StatusOK, resp)
}

// GetDashboardCharts godoc
// @Summary Get dashboard charts data
// @Description Get cash flow trends and category breakdown charts using Hybrid Logic
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param range query string false "Date range" Enums(7d, 1m, 3m)
// @Success 200 {object} map[string]interface{} "Dashboard charts retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/dashboard/charts [get]
func (c *Controller) GetDashboardCharts(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	timeRange := TimeRange(ctx.DefaultQuery("range", ""))
	if timeRange != "" && timeRange != TimeRange7Days && timeRange != TimeRange1Month && timeRange != TimeRange3Month {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid time_range, allowed values: 7d, 1m, 3m")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	charts, err := c.service.GetDashboardCharts(ctx, userID.(string), timeRange)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Dashboard charts retrieved successfully", charts)
	ctx.JSON(http.StatusOK, resp)
}

// SyncDailySummaries godoc
// @Summary Sync daily summaries
// @Description Sync daily summaries for all users starting from a specific date. Delete first.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Daily summaries synced successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/dashboard/sync [post]
func (c *Controller) SyncDailySummaries(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	if startDateStr == "" {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "start_date is required")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	err = c.dailySummaryService.SyncDailySummaries(ctx, startDate)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Daily summaries synced successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
