package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log-snare/web/service"
)

type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler initializes a new user handler with the given user service
func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{UserService: us}
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

	c.Redirect(302, "/dashboard")
}
