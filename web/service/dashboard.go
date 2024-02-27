package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log-snare/web/data"
)

type DashboardService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewDashboardService(db *gorm.DB, logger *zap.Logger) *DashboardService {
	return &DashboardService{DB: db, Logger: logger}
}

type DashboardSummaryCount struct {
	EmployeeCount           int64
	EmployeeHighSalaryCount int64
	UserCount               int64
	AdminCount              int64
}

func (s *DashboardService) GetSummaryCounts(companyId int) (retval DashboardSummaryCount, err error) {

	s.DB.Model(&data.Employee{}).Where("company_id = ?", companyId).Count(&retval.EmployeeCount)
	s.DB.Model(&data.Employee{}).Where("company_id = ? AND salary >= ?", companyId, 100000).Count(&retval.EmployeeHighSalaryCount)

	s.DB.Model(&data.User{}).Where("company_id = ?", companyId).Count(&retval.UserCount)
	s.DB.Model(&data.User{}).Where("company_id = ? AND role = 1", companyId).Count(&retval.AdminCount)

	return retval, nil
}

type CompletedChallenges struct {
	OneComplete   bool
	TwoComplete   bool
	ThreeComplete bool
}

func (s *DashboardService) GetCompletedChallenges() (retval CompletedChallenges, err error) {
	var settingValue data.SettingValue
	err = s.DB.Model(&data.SettingValue{}).Where("key = ? AND value = ?", 1, 1).First(&settingValue).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return retval, fmt.Errorf("unable to retrieve challenge info: %w", err.Error())
		} else {
			retval.OneComplete = false
		}
	} else {
		retval.OneComplete = true
	}

	return retval, nil
}
