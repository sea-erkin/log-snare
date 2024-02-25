package service

import (
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
