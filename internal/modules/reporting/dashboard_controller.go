package reporting

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardController struct {
	service *Service
}

func NewDashboardController(service *Service) *DashboardController {
	return &DashboardController{service: service}
}

// ============================================================================
// KPI ENDPOINTS
// ============================================================================

// GetKPIs returns dashboard KPI cards
// GET /api/v1/dashboard/kpis
func (c *DashboardController) GetKPIs(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	kpis, err := c.service.GetDashboardKPIs(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kpis,
	})
}

// ============================================================================
// CHART ENDPOINTS
// ============================================================================

// GetCharts returns dashboard charts
// GET /api/v1/dashboard/charts
func (c *DashboardController) GetCharts(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	charts, err := c.service.GetDashboardCharts(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    charts,
	})
}

// ============================================================================
// DAILY REPORTS WITH DATE RANGE FILTER
// ============================================================================

// GetDailyReportsByDateRange returns daily reports for a date range
// GET /api/v1/dashboard/reports/daily?start_date=2024-01-01&end_date=2024-01-31
// Query Parameters:
//   - start_date: YYYY-MM-DD (required)
//   - end_date: YYYY-MM-DD (required)
func (c *DashboardController) GetDailyReportsByDateRange(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Parse date range
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (use YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (use YYYY-MM-DD)"})
		return
	}

	// Ensure start_date <= end_date
	if startDate.After(endDate) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before or equal to end_date"})
		return
	}

	// Fetch reports
	reports, err := c.service.repo.GetDailyReportsByDateRange(ctx, objID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"start_date": startDateStr,
			"end_date":   endDateStr,
			"count":      len(reports),
			"reports":    reports,
		},
	})
}

// ============================================================================
// MONTHLY REPORTS WITH DATE RANGE FILTER
// ============================================================================

// GetMonthlyReportsByDateRange returns monthly reports for a date range
// GET /api/v1/dashboard/reports/monthly?start_month=2024-01&end_month=2024-12
// Query Parameters:
//   - start_month: YYYY-MM (required)
//   - end_month: YYYY-MM (required)
func (c *DashboardController) GetMonthlyReportsByDateRange(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Parse month range
	startMonthStr := ctx.Query("start_month")
	endMonthStr := ctx.Query("end_month")

	if startMonthStr == "" || endMonthStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "start_month and end_month are required",
		})
		return
	}

	startMonth, err := time.Parse("2006-01", startMonthStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_month format (use YYYY-MM)"})
		return
	}

	endMonth, err := time.Parse("2006-01", endMonthStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_month format (use YYYY-MM)"})
		return
	}

	// Ensure start_month <= end_month
	if startMonth.After(endMonth) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_month must be before or equal to end_month"})
		return
	}

	// Fetch summaries
	summaries, err := c.service.repo.GetMonthlySummariesByRange(ctx, objID, startMonth, endMonth)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"start_month": startMonthStr,
			"end_month":   endMonthStr,
			"count":       len(summaries),
			"summaries":   summaries,
		},
	})
}

// ============================================================================
// INCOME/EXPENSE CHART WITH DATE RANGE FILTER
// ============================================================================

// GetIncomeExpenseChart returns monthly income vs expense chart
// GET /api/v1/dashboard/charts/income-expense?months=12
// Query Parameters:
//   - months: number of months to include (default: 12, max: 60)
func (c *DashboardController) GetIncomeExpenseChart(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Parse months parameter
	monthsStr := ctx.DefaultQuery("months", "12")
	months := 12
	if monthsStr != "" {
		if m, err := parseIntParam(monthsStr); err == nil {
			if m > 0 && m <= 60 {
				months = m
			}
		}
	}

	chart, err := c.service.aggregationHelper.GetMonthlyIncomeExpenseChart(ctx, objID, months)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"months": months,
			"chart":  chart,
		},
	})
}

// ============================================================================
// CATEGORY DISTRIBUTION WITH DATE RANGE FILTER
// ============================================================================

// GetCategoryDistribution returns expense distribution by category
// GET /api/v1/dashboard/charts/category-distribution?month=2024-01
// Query Parameters:
//   - month: YYYY-MM (optional, defaults to current month)
func (c *DashboardController) GetCategoryDistribution(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Parse month parameter
	monthStr := ctx.Query("month")
	month := time.Now()

	if monthStr != "" {
		if m, err := time.Parse("2006-01", monthStr); err == nil {
			month = m
		}
	}

	// Normalize to first day of month
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)

	distribution, err := c.service.aggregationHelper.GetExpenseCategoryDistribution(ctx, objID, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"month":        month.Format("2006-01"),
			"distribution": distribution,
		},
	})
}

// ============================================================================
// POCKET DISTRIBUTION
// ============================================================================

// GetPocketDistribution returns balance distribution by pocket type
// GET /api/v1/dashboard/charts/pocket-distribution
func (c *DashboardController) GetPocketDistribution(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	distribution, err := c.service.aggregationHelper.GetPocketBalanceDistribution(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    distribution,
	})
}

// ============================================================================
// SUMMARY ENDPOINT
// ============================================================================

// GetDashboardSummary returns complete dashboard summary with all data
// GET /api/v1/dashboard/summary?start_date=2024-01-01&end_date=2024-01-31
// Query Parameters:
//   - start_date: YYYY-MM-DD (optional)
//   - end_date: YYYY-MM-DD (optional)
func (c *DashboardController) GetDashboardSummary(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Get KPIs
	kpis, err := c.service.GetDashboardKPIs(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get charts
	charts, err := c.service.GetDashboardCharts(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get date range if provided
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var dateRangeReports interface{}
	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil && !startDate.After(endDate) {
				reports, err := c.service.repo.GetDailyReportsByDateRange(ctx, objID, startDate, endDate)
				if err == nil {
					dateRangeReports = reports
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"kpis":               kpis,
			"charts":             charts,
			"date_range_reports": dateRangeReports,
		},
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func parseIntParam(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
