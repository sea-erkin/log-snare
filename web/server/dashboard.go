package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log-snare/web/service"
)

type DashboardHandler struct {
	DashboardService *service.DashboardService
}

// NewDashboardHandler initializes a new user handler with the given user service
func NewDashboardHandler(us *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{DashboardService: us}
}

func (h *DashboardHandler) SummaryCounts(c *gin.Context) {

	companyId := 1

	summaryCounts, err := h.DashboardService.GetSummaryCounts(companyId)
	if err != nil {
		h.DashboardService.Logger.Error("unable get summary counts", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	c.HTML(200, "dashboard.html", gin.H{
		"CurrentRoute":    "/dashboard",
		"EmployeeCount":   summaryCounts.EmployeeCount,
		"AdminCount":      summaryCounts.AdminCount,
		"UserCount":       summaryCounts.UserCount,
		"HighSalaryCount": summaryCounts.EmployeeHighSalaryCount,
	})
}
