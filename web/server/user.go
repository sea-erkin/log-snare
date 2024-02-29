package server

import (
	"encoding/base64"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"log-snare/web/data"
	"log-snare/web/service"
	"strconv"
	"strings"
)

type UserHandler struct {
	UserService    *service.UserService
	SettingService *service.SettingsService
}

// NewUserHandler initializes a new user handler with the given user service
func NewUserHandler(us *service.UserService, ss *service.SettingsService) *UserHandler {
	return &UserHandler{UserService: us, SettingService: ss}
}

func (uh *UserHandler) Login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	userExists, err := uh.UserService.CheckUsernameExists(uh.UserService.DB, username)
	if err != nil {
		uh.UserService.Logger.Error("unable to check if user exists", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	if !userExists {
		uh.UserService.Logger.Info("attempted user does not exist",
			zap.String("username", username),
			zap.String("eventType", "security"),
			zap.String("securityType", "bruteforce-possible"),
			zap.String("eventCategory", "logon"),
			zap.String("clientIp", c.ClientIP()),
			zap.Bool("success", false),
		)
		c.HTML(200, "login.html", gin.H{
			"Error": "Login failed, attempt has been logged.",
		})
		return
	}

	// if we get any error, log as login failure
	user, err := uh.UserService.GetUserByUsernameAndPassword(username, password)
	if err != nil {
		uh.UserService.Logger.Info("logon failure",
			zap.String("username", username),
			zap.String("eventType", "security"),
			zap.String("securityType", "bruteforce-possible"),
			zap.String("eventCategory", "logon"),
			zap.Bool("success", false),
			zap.String("clientIp", c.ClientIP()),
			zap.Error(err),
		)
		c.HTML(200, "login.html", gin.H{
			"Error": "Login failed, attempt has been logged.",
		})
		return
	}

	// if we got this far, we've logged in
	uh.UserService.Logger.Info("logon success",
		zap.String("username", username),
		zap.String("eventType", "security"),
		zap.String("eventCategory", "logon"),
		zap.Bool("success", true),
		zap.String("clientIp", c.ClientIP()),
	)

	session := sessions.Default(c)
	session.Set("user", user)
	err = session.Save()
	if err != nil {
		uh.UserService.Logger.Error("failed to save session", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	c.Redirect(302, "/app/docs")
}

func (h *UserHandler) EnableAdmin(c *gin.Context) {

	session := sessions.Default(c)
	userSession := session.Get("user").(data.UserSafe)

	if data.ValidationEnabled() {
		if userSession.Role != service.RoleAdmin {
			h.UserService.Logger.Warn("user is trying to enable admin, but they are a basic user",
				zap.String("username", userSession.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}
	}

	// Update user role to Admin
	err := h.UserService.DB.Model(&data.User{}).Where("id = ?", userSession.UserId).Update("role", 1).Error
	if err != nil {
		h.UserService.Logger.Error("unable to set user admin", zap.Error(err))
	}

	// Challenge two complete, if the user was not an admin
	if userSession.Role == 2 {
		h.SettingService.ChallengeComplete("2")
	}

	// set in session
	userSession.Role = 1
	session.Set("user", userSession)
	err = session.Save()
	if err != nil {
		h.UserService.Logger.Error("failed to save session", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *UserHandler) DisableAdmin(c *gin.Context) {

	session := sessions.Default(c)
	userSession := session.Get("user").(data.UserSafe)

	if data.ValidationEnabled() {
		if userSession.Role != service.RoleAdmin {
			h.UserService.Logger.Warn("user is trying to disable admin, but they are a basic user",
				zap.String("username", userSession.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}
	}

	// set in DB
	err := h.UserService.DB.Model(&data.User{}).Where("id = ?", userSession.UserId).Update("role", 2).Error
	if err != nil {
		h.UserService.Logger.Error("unable to set user admin", zap.Error(err))
	}

	// set in session
	userSession.Role = 2
	session.Set("user", userSession)
	err = session.Save()
	if err != nil {
		h.UserService.Logger.Error("failed to save session", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *UserHandler) Users(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("user").(data.UserSafe)

	// expecting input as base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(c.Param("id"))
	if err != nil {
		h.UserService.Logger.Warn("user has provided invalid base64 when expecting base64",
			zap.String("username", user.Username),
			zap.String("eventType", "security"),
			zap.String("securityType", "tamper-certain"),
			zap.String("eventCategory", "validation"),
			zap.String("clientIp", c.ClientIP()),
		)
		c.HTML(404, "error-404.html", nil)
		return
	}

	id := strings.TrimPrefix(string(decodedBytes), "CompanyId:")
	intId, err := strconv.Atoi(id)
	if err != nil {
		h.UserService.Logger.Warn("user has provided an invalid CompanyId when expecting an integer",
			zap.String("username", user.Username),
			zap.String("eventType", "security"),
			zap.String("securityType", "tamper-certain"),
			zap.String("eventCategory", "validation"),
			zap.String("clientIp", c.ClientIP()),
		)
		c.HTML(404, "error-404.html", nil)
		return
	}

	if data.ValidationEnabled() {
		if user.Role == service.RoleUser {
			h.UserService.Logger.Warn("user is not an admin, but trying to access manage users interface",
				zap.String("username", user.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}

		if user.CompanyId != intId {
			h.UserService.Logger.Warn("user is trying to manage users for a company Id they do not belong too.",
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

	users := h.UserService.UsersByCompanyId(intId)
	var compName string
	if len(users) > 0 {
		compName = users[0].CompanyName
	}

	c.HTML(200, "users.html", gin.H{
		"CurrentRoute":        "/users",
		"Users":               users,
		"UserCompanyName":     compName,
		"UserCount":           len(users),
		"UserCompanyId":       user.CompanyId,
		"ManageUserCompanyId": manageUserCompanyIdentifier(user.CompanyId),

		// common data can be moved to middleware
		"CompanyName":       user.CompanyName,
		"UserInitial":       string(strings.ToUpper(user.Username)[0]),
		"UserRole":          user.Role,
		"ValidationEnabled": data.ValidationEnabled(),
	})

}

func (h *UserHandler) Impersonate(c *gin.Context) {

	session := sessions.Default(c)
	sessionUser := session.Get("user").(data.UserSafe)

	identifier, err := ksuid.Parse(c.Param("id"))
	if err != nil {
		h.UserService.Logger.Warn("unable to parse KSUID identifier for manage users",
			zap.String("username", sessionUser.Username),
			zap.String("eventType", "security"),
			zap.String("securityType", "tamper-certain"),
			zap.String("eventCategory", "validation"),
			zap.String("clientIp", c.ClientIP()),
		)
		c.HTML(404, "error-404.html", nil)
		return
	}

	// set user
	targetUser := h.UserService.UserByIdentifier(identifier.String())
	session.Set("user", targetUser)
	err = session.Save()
	if err != nil {
		h.UserService.Logger.Error("failed to save session", zap.Error(err))
		c.Redirect(500, "/")
		return
	}

	if data.ValidationEnabled() {
		if sessionUser.Role == service.RoleUser {
			h.UserService.Logger.Warn("user is not an admin, but is trying to impersonate users",
				zap.String("username", sessionUser.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}

		if sessionUser.CompanyId != targetUser.CompanyId {
			h.UserService.Logger.Warn("user is trying to impersonate someone outside of their company",
				zap.String("username", sessionUser.Username),
				zap.String("eventType", "security"),
				zap.String("securityType", "tamper-certain"),
				zap.String("eventCategory", "validation"),
				zap.String("clientIp", c.ClientIP()),
			)
			c.HTML(404, "error-404.html", nil)
			return
		}
	}

	// challenge complete if we made it this far
	if sessionUser.CompanyId != targetUser.CompanyId && targetUser.Role == service.RoleAdmin {
		h.SettingService.ChallengeComplete("3")
		c.JSON(200, gin.H{
			"success": false,
			"message": "Congrats admin of not your company.",
		})
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
