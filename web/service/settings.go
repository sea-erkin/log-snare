package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log-snare/web/data"
)

type SettingsService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewSettingsService(db *gorm.DB, logger *zap.Logger) *SettingsService {
	return &SettingsService{DB: db, Logger: logger}
}

func (s *SettingsService) ChallengeComplete(challenge string) {
	setting := data.SettingValue{}
	s.DB.Model(&data.SettingValue{}).Where("key = ?", "1").First(&setting)
	setting.Value = true
	s.DB.Save(setting)
}
