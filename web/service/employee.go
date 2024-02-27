package service

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log-snare/web/data"
	"math"
	"strings"
)

type EmployeeService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewEmployeeService(db *gorm.DB, logger *zap.Logger) *EmployeeService {
	return &EmployeeService{DB: db, Logger: logger}
}

func (s *EmployeeService) EmployeesByCompanyId(companyId int) (retval []data.Employee) {
	s.DB.Preload("Company").Where("company_id = ?", companyId).Find(&retval)

	for i, employee := range retval {
		roundedNumber := math.Round(employee.Salary*100) / 100
		formattedNumber := fmt.Sprintf("$%v", roundedNumber)
		retval[i].DisplaySalary = formattedNumber
		retval[i].DisplayDOB = employee.DateOfBirth.Format("01/02/2006")
		retval[i].DisplaySSN = "###-##-" + strings.Split(employee.SSN, "-")[2]
	}

	return retval
}
