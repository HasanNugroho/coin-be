package reporting

import (
<<<<<<< Updated upstream
	"fmt"
	"net/http"
	"time"

=======
	"net/http"
	"strconv"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
>>>>>>> Stashed changes
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

<<<<<<< Updated upstream
type DashboardController struct {
	service *DashboardService
}

func NewDashboardController(service *DashboardService) *DashboardController {
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
=======
type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// GetDailyReport godoc
// @Summary Get daily financial report
// @Description Get precomputed daily financial report for a specific date
// @Tags Reports
// @Produce json
// @Param date query string true "Report date (YYYY-MM-DD)"
// @Param include_details query boolean false "Include detailed breakdowns"
// @Security BearerAuth
// @Success 200 {object} DailyFinancialReport
// @Success 202 {object} map[string]string "Report queued for generation"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/reports/daily [get]
func (c *Controller) GetDailyReport(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "date parameter required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD"))
		return
	}

	userOID, ok := userID.(primitive.ObjectID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	report, err := c.service.GetDailyReport(ctx, userOID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	if report == nil {
		ctx.JSON(http.StatusAccepted, utils.NewResponse(http.StatusAccepted, "report queued for generation", nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("daily report retrieved", report))
}

// GenerateDailyReport godoc
// @Summary Generate daily financial report
// @Description Trigger generation of daily financial report for a specific date
// @Tags Reports
// @Produce json
// @Param date query string true "Report date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} DailyFinancialReport
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/reports/daily/generate [post]
func (c *Controller) GenerateDailyReport(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "date parameter required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	report, err := c.service.GenerateDailyReport(ctx, userOID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("daily report generated", report))
}

// GetDashboardKPIs godoc
// @Summary Get dashboard KPI cards
// @Description Get key performance indicators for dashboard
// @Tags Dashboard
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/kpis [get]
func (c *Controller) GetDashboardKPIs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	month := ctx.Query("month")
	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	// Get current snapshot for total balance
	snapshot, err := c.service.repo.GetDailySnapshot(ctx, userOID, time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// Get monthly summary
	summary, err := c.service.GetMonthlySummary(ctx, userOID, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	kpis := gin.H{
		"total_balance":               0.0,
		"total_income_current_month":  0.0,
		"total_expense_current_month": 0.0,
		"free_money_total":            0.0,
		"net_change_current_month":    0.0,
	}

	if snapshot != nil {
		kpis["total_balance"] = snapshot.TotalBalance
		kpis["free_money_total"] = snapshot.FreeMoneyTotal
	}

	if summary != nil {
		kpis["total_income_current_month"] = summary.TotalIncome
		kpis["total_expense_current_month"] = summary.TotalExpense
		kpis["net_change_current_month"] = summary.TotalIncome - summary.TotalExpense
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("dashboard KPIs retrieved", kpis))
}

// GetMonthlyTrendChart godoc
// @Summary Get monthly income vs expense trend
// @Description Get 12-month trend of income vs expense
// @Tags Dashboard
// @Produce json
// @Param months query int false "Number of months (default 12, max 36)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/charts/monthly-trend [get]
func (c *Controller) GetMonthlyTrendChart(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	monthsStr := ctx.Query("months")
	months := int64(12)
	if monthsStr != "" {
		if m, err := strconv.ParseInt(monthsStr, 10, 64); err == nil {
			months = m
		}
	}

	summaries, err := c.service.GetMonthlyTrend(ctx, userOID, months)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	type TrendData struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Net     float64 `json:"net"`
	}

	data := []TrendData{}
	for _, summary := range summaries {
		data = append(data, TrendData{
			Month:   summary.YearMonth,
			Income:  summary.TotalIncome,
			Expense: summary.TotalExpense,
			Net:     summary.TotalIncome - summary.TotalExpense,
		})
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("monthly trend retrieved", gin.H{"data": data}))
}

// GetPocketDistributionChart godoc
// @Summary Get balance distribution per pocket
// @Description Get pie chart data for pocket balance distribution
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/charts/pocket-distribution [get]
func (c *Controller) GetPocketDistributionChart(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	snapshot, err := c.service.repo.GetDailySnapshot(ctx, userOID, time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	type PocketData struct {
		PocketID   string  `json:"pocket_id"`
		PocketName string  `json:"pocket_name"`
		PocketType string  `json:"pocket_type"`
		Balance    float64 `json:"balance"`
		Percentage float64 `json:"percentage"`
	}

	data := []PocketData{}
	if snapshot != nil {
		for _, pocket := range snapshot.PocketBalances {
			percentage := 0.0
			if snapshot.TotalBalance > 0 {
				percentage = (pocket.Balance / snapshot.TotalBalance) * 100
			}
			data = append(data, PocketData{
				PocketID:   pocket.PocketID.Hex(),
				PocketName: pocket.PocketName,
				PocketType: pocket.PocketType,
				Balance:    pocket.Balance,
				Percentage: percentage,
			})
		}
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("pocket distribution retrieved", gin.H{"data": data}))
}

// GetExpenseByCategoryChart godoc
// @Summary Get expense distribution by category
// @Description Get pie chart data for expense distribution
// @Tags Dashboard
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/charts/expense-by-category [get]
func (c *Controller) GetExpenseByCategoryChart(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	month := ctx.Query("month")
	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	summary, err := c.service.GetMonthlySummary(ctx, userOID, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	type CategoryData struct {
		CategoryID       string  `json:"category_id"`
		CategoryName     string  `json:"category_name"`
		Amount           float64 `json:"amount"`
		Percentage       float64 `json:"percentage"`
		TransactionCount int32   `json:"transaction_count"`
	}

	data := []CategoryData{}
	if summary != nil {
		for _, cat := range summary.ExpenseByCategory {
			percentage := 0.0
			if summary.TotalExpense > 0 {
				percentage = (cat.Amount / summary.TotalExpense) * 100
			}
			data = append(data, CategoryData{
				CategoryID:       cat.CategoryID.Hex(),
				CategoryName:     cat.CategoryName,
				Amount:           cat.Amount,
				Percentage:       percentage,
				TransactionCount: cat.TransactionCount,
			})
		}
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("expense by category retrieved", gin.H{"data": data}))
}

// GetAIFinancialContext godoc
// @Summary Get AI-ready financial context
// @Description Get precomputed financial context for AI chatbot
// @Tags AI
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AIFinancialContext
// @Failure 401 {object} map[string]string
// @Router /api/v1/ai/financial-context [get]
func (c *Controller) GetAIFinancialContext(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, ok := userID.(primitive.ObjectID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	aiContext, err := c.service.GetAIFinancialContext(ctx, userOID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	if aiContext == nil {
		ctx.JSON(http.StatusAccepted, utils.NewResponse(http.StatusAccepted, "context being generated", nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("AI financial context retrieved", aiContext))
}

// HealthCheck godoc
// @Summary Health check for reporting service
// @Description Check if reporting service is healthy
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health/reporting [get]
func (c *Controller) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("reporting service healthy", gin.H{"status": "healthy", "service": "reporting"}))
}

// GetRealtimeDashboardKPIs godoc
// @Summary Get real-time dashboard KPI cards
// @Description Get real-time key performance indicators by querying live collections
// @Tags Dashboard
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/realtime/kpis [get]
func (c *Controller) GetRealtimeDashboardKPIs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	month := ctx.Query("month")

	kpis, err := c.service.GetRealtimeDashboardKPIs(ctx, userOID, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("real-time dashboard KPIs retrieved", kpis))
}

// GetRealtimePocketDistributionChart godoc
// @Summary Get real-time balance distribution per pocket
// @Description Get real-time pie chart data for pocket balance distribution
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/realtime/charts/pocket-distribution [get]
func (c *Controller) GetRealtimePocketDistributionChart(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	data, err := c.service.GetRealtimePocketDistribution(ctx, userOID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("real-time pocket distribution retrieved", gin.H{"data": data}))
}

// GetRealtimeExpenseByCategoryChart godoc
// @Summary Get real-time expense distribution by category
// @Description Get real-time pie chart data for expense distribution
// @Tags Dashboard
// @Produce json
// @Param month query string false "Month (YYYY-MM), defaults to current"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/dashboard/realtime/charts/expense-by-category [get]
func (c *Controller) GetRealtimeExpenseByCategoryChart(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized"))
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "invalid user id"))
		return
	}

	month := ctx.Query("month")

	data, err := c.service.GetRealtimeExpenseByCategory(ctx, userOID, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("real-time expense by category retrieved", gin.H{"data": data}))
>>>>>>> Stashed changes
}
