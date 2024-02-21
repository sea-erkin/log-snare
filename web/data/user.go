package data

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name  string
	Users []User `gorm:"foreignKey:CompanyId"` // Indicates a one-to-many relationship
}

type User struct {
	gorm.Model
	CompanyId int    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Username  string `gorm:"uniqueIndex"`
	Password  string
	Role      int
	Active    bool `gorm:"default:true"`
}

// SetPassword hashes the password and stores it in the User struct
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
