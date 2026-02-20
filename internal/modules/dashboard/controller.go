package dashboard

import (
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
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
