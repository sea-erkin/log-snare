package data

import (
	"crypto/rand"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type Company struct {
	gorm.Model
	Name       string
	Identifier string
	Users      []User     `gorm:"foreignKey:CompanyId"`
	Employees  []Employee `gorm:"foreignKey:CompanyId"`
}

func (m *Company) BeforeCreate(tx *gorm.DB) (err error) {
	m.Identifier = ksuid.New().String()
	return
}

type User struct {
	gorm.Model
	Identifier string `gorm:"index;unique"` // KSUID identifier for demo
	CompanyId  int    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Username   string `gorm:"uniqueIndex"`
	Password   string
	Role       int
	Active     bool `gorm:"default:true"`
	Company    Company
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.Identifier = ksuid.New().String()
	return
}

// UserSafe does not contain a password
type UserSafe struct {
	CompanyId   int
	Username    string
	Role        int
	Active      bool
	CompanyName string
	UserId      uint
	Identifier  string
	Created     string
}

func (m *User) ToUserSafe() UserSafe {
	return UserSafe{
		CompanyId:   m.CompanyId,
		Username:    m.Username,
		Role:        m.Role,
		Active:      m.Active,
		UserId:      m.ID,
		CompanyName: m.Company.Name,
		Identifier:  m.Identifier,
		Created:     m.CreatedAt.Format("01/02/2006"),
	}
}

// SetPassword hashes the password and stores it in the User struct
func (m *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hash)
	return nil
}

func CreateUserWithPassword(username string, role int, logPassword bool) User {
	user := User{
		Username: username,
		Role:     role,
	}
	password := generatePassword(12)
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
		"0123456789"

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
