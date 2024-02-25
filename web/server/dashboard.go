package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log-snare/web/data"
	"log-snare/web/service"
	"strings"
)

type DashboardHandler struct {
	DashboardService *service.DashboardService
}

// NewDashboardHandler initializes a new user handler with the given user service
func NewDashboardHandler(us *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{DashboardService: us}
}

func (h *DashboardHandler) SummaryCounts(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("user").(data.UserSafe)

	summaryCounts, err := h.DashboardService.GetSummaryCounts(user.CompanyId)
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

		// common data can be moved to middleware
		"CompanyName": user.CompanyName,
		"UserInitial": string(strings.ToUpper(user.Username)[0]),
		"UserRole":    user.Role,
	})
}
