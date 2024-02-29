package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log-snare/web/data"
	"log-snare/web/service"
	"strconv"
	"strings"
)

type EmployeeHandler struct {
	EmployeeService *service.EmployeeService
	SettingsService *service.SettingsService
	UserService     *service.UserService
}

// NewEmployeeHandler initializes a new user handler with the given user service
func NewEmployeeHandler(us *service.EmployeeService, ss *service.SettingsService, uss *service.UserService) *EmployeeHandler {
	return &EmployeeHandler{
		EmployeeService: us,
		SettingsService: ss,
		UserService:     uss,
	}
}

func (h *EmployeeHandler) Employees(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("user").(data.UserSafe)

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		h.EmployeeService.Logger.Warn("unable to parse provided ID as an integer",
			zap.String("username", user.Username),
			zap.String("eventType", "security"),
			zap.String("securityType", "tamper-possible"),
			zap.String("eventCategory", "validation"),
			zap.String("clientIp", c.ClientIP()),
		)
		c.HTML(404, "error-404.html", nil)
		return
	}

	if data.ValidationEnabled() {
		if user.CompanyId != intId {
			h.EmployeeService.Logger.Warn("user is trying to access a company ID that is not theirs",
				zap.String("username", user.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}
	}

	employees := h.EmployeeService.EmployeesByCompanyId(intId)
	var employeeCompanyName string
	if len(employees) > 0 {
		employeeCompanyName = employees[0].Company.Name
	}

	// if we made it this far, and we're serving content for another company ID, challenge one is completed
	if user.CompanyId != intId {
		h.SettingsService.ChallengeComplete("1")
	}

	c.HTML(200, "employees.html", gin.H{
		"CurrentRoute":        "/employees",
		"Employees":           employees,
		"EmployeeCompanyName": employeeCompanyName,
		"EmployeeCount":       len(employees),
		"UserCompanyId":       user.CompanyId,
		"ManageUserCompanyId": manageUserCompanyIdentifier(user.CompanyId),

		// common data can be moved to middleware
		"CompanyName":       user.CompanyName,
		"UserInitial":       string(strings.ToUpper(user.Username)[0]),
		"UserRole":          user.Role,
		"ValidationEnabled": data.ValidationEnabled(),
	})

}
