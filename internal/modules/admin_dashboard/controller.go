package admin_dashboard

import (
	"net/http"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// GetAdminSummary godoc
// @Summary Get admin dashboard summary
// @Description Get global statistics for administrators with date filters
// @Tags Admin
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Admin dashboard summary retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Security BearerAuth
// @Router /v1/admin/dashboard [get]
func (c *Controller) GetAdminSummary(ctx *gin.Context) {
	startDateStr := ctx.DefaultQuery("start_date", "")
	endDateStr := ctx.DefaultQuery("end_date", "")

	startDate := time.Now().AddDate(0, 0, -30) // Default 30 days
	endDate := time.Now()

	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	}

	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = t
		}
	}

	summary, err := c.service.GetAdminDashboardSummary(ctx, startDate, endDate)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Admin dashboard summary retrieved successfully", summary)
	ctx.JSON(http.StatusOK, resp)
}
