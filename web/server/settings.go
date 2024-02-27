package server

import (
	"github.com/gin-gonic/gin"
	"log-snare/web/data"
	"log-snare/web/service"
)

type SettingsHandler struct {
	SettingsService *service.SettingsService
}

// NewSettingsHandler initializes a new user handler with the given user service
func NewSettingsHandler(us *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{SettingsService: us}
}

func (h *SettingsHandler) EnableValidation(c *gin.Context) {
	data.SetValidation(true)
	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *SettingsHandler) DisableValidation(c *gin.Context) {
	data.SetValidation(false)
	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *SettingsHandler) EnableLogging(c *gin.Context) {
	data.SetLogging(true)
	c.JSON(200, gin.H{
		"success": true,
	})
}

func (h *SettingsHandler) DisableLogging(c *gin.Context) {
	data.SetLogging(false)
	c.JSON(200, gin.H{
		"success": true,
	})
}
