package data

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type Company struct {
	gorm.Model
	Name      string
	Users     []User     `gorm:"foreignKey:CompanyId"`
	Employees []Employee `gorm:"foreignKey:CompanyId"`
}

type User struct {
	gorm.Model
	CompanyId int    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Username  string `gorm:"uniqueIndex"`
	Password  string
	Role      int
	Active    bool `gorm:"default:true"`
}

// UserSafe does not contain a password
type UserSafe struct {
	CompanyId int
	Username  string
	Role      int
	Active    bool
}

func (u *User) UserToUserSafe() UserSafe {
	return UserSafe{
		CompanyId: u.CompanyId,
		Username:  u.Username,
		Role:      u.Role,
		Active:    u.Active,
	}
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

func CreateUserWithPassword(username string, role int, logPassword bool) User {
	user := User{
		Username: username,
		Role:     role,
	}
	password := generatePassword(24)
	if logPassword {
		log.Printf("[DATA] %s password: %s\n", username, password)
	}

	err := user.SetPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*"

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = charset[b[i]%byte(len(charset))]
	}

	return string(password)
}
