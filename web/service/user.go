package service

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log-snare/web/data"
)

type UserService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewUserService(db *gorm.DB, logger *zap.Logger) *UserService {
	return &UserService{DB: db, Logger: logger}
}

func (us *UserService) CheckUsernameExists(db *gorm.DB, username string) (bool, error) {
	var user data.User
	result := db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (us *UserService) GetUserByUsernameAndPassword(username, password string) (retval data.UserSafe, err error) {
	var user data.User

	res := us.DB.Preload("Company").Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return retval, errors.New("username not found")
		}
		return retval, res.Error
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return retval, err
	}

	return user.ToUserSafe(), nil
}

func (us *UserService) UsersByCompanyId(companyId int) (retval []data.UserSafe) {
	var users []data.User
	us.DB.Preload("Company").Where("company_id = ?", companyId).Find(&users)

	for _, user := range users {
		retval = append(retval, user.ToUserSafe())
	}

	return retval
}

func (us *UserService) UserByIdentifier(identifier string) (retval data.UserSafe) {
	var user data.User
	us.DB.Preload("Company").Where("identifier = ?", identifier).First(&user)
	return user.ToUserSafe()
}
