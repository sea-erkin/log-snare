package data

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

// Employee represents an HR employee with various personal and professional details.
type Employee struct {
	gorm.Model            // Embedding gorm.Model adds fields ID, CreatedAt, UpdatedAt, DeletedAt
	CompanyId     int     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CompanyKSUID  int     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	EmployeeKSUID string  `gorm:"index;unique"` // KSUID identifier for demo
	FirstName     string  `gorm:"type:varchar(100);not null"`
	LastName      string  `gorm:"type:varchar(100);not null"`
	Email         string  `gorm:"type:varchar(100);not null"`
	SSN           string  `gorm:"type:varchar(11);unique;not null"` // Assuming US format XXX-XX-XXXX
	Salary        float64 `gorm:"not null"`
	DateOfBirth   time.Time
	Company       Company
	DisplaySalary string `gorm:"-"`
	DisplayDOB    string `gorm:"-"`
	DisplaySSN    string `gorm:"-"`
}

func (m *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	m.EmployeeKSUID = ksuid.New().String()
	return
}

func GenerateEmployee(emailSuffix string) Employee {

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	user := Employee{}
	user.FirstName, user.LastName = getRandomNames(r)
	user.Email = strings.ToLower(user.FirstName) + "." + strings.ToLower(user.LastName) + "@" + emailSuffix
	user.SSN = generateRandomSSN(r)
	user.Salary = generateRandomSalary(r)
	user.DateOfBirth = generateRandomDOB(r)

	return user
}

func getRandomNames(r *rand.Rand) (string, string) {
	index := r.Intn(len(firstNames) - 1)
	indexLast := r.Intn(len(lastNames) - 1)
	return firstNames[index], lastNames[indexLast]
}

func generateRandomSSN(r *rand.Rand) string {
	area := r.Intn(900) + 100 // Avoid generating leading zeros, range: 100-999
	group := r.Intn(100)      // Range: 00-99
	serial := r.Intn(10000)   // Range: 0000-9999

	ssn := fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
	return ssn
}

func generateRandomSalary(r *rand.Rand) float64 {
	min := 40000.0  // Minimum salary
	max := 125000.0 // Maximum salary

	// Generate a random float64 between min and max
	salary := min + r.Float64()*(max-min)
	return salary
}

func generateRandomDOB(r *rand.Rand) time.Time {
	now := time.Now()
	// Calculate 100 years ago from now
	hundredYearsAgo := now.AddDate(-100, 0, 0)
	// Calculate the difference in days between now and 100 years ago
	daysDiff := now.Sub(hundredYearsAgo).Hours() / 24
	// Generate a random number of days within that range
	randomDays := r.Intn(int(daysDiff))
	// Calculate the random date of birth by adding the random number of days to 100 years ago
	randomDOB := hundredYearsAgo.AddDate(0, 0, randomDays)
	return randomDOB
}

var firstNames = []string{
	"Emma", "Liam", "Olivia", "Noah", "Ava", "Oliver", "Isabella", "Mason", "Sophia", "Logan",
	"Mia", "Lucas", "Charlotte", "Ethan", "Amelia", "Elijah", "Harper", "Benjamin", "Evelyn", "Sebastian",
	"Abigail", "Jackson", "Emily", "Aiden", "Elizabeth", "Matthew", "Mila", "Samuel", "Ella", "David",
	"Scarlett", "Joseph", "Madison", "Carter", "Layla", "Owen", "Chloe", "Wyatt", "Grace", "John",
	"Ellie", "Jack", "Zoey", "Luke", "Penelope", "Jayden", "Riley", "Dylan", "Nora", "Leo",
	"Lily", "Alexander", "Hannah", "Grayson", "Luna", "Michael", "Zoe", "James", "Stella", "Ezra",
	"Addison", "Isaac", "Lillian", "Gabriel", "Aubrey", "Julian", "Audrey", "Mateo", "Elliot", "Ian",
	"Rose", "Josiah", "Violet", "Theodore", "Claire", "Avery", "Lincoln", "Lucy", "Asher", "Caroline",
	"John", "Nova", "Jonathan", "Genesis", "Xavier", "Emilia", "Jaxon", "Kennedy", "Isaiah", "Samantha",
	"Elias", "Maya", "Aaron", "Willow", "Charles", "Kinsley", "Christopher", "Naomi", "Cameron", "Aaliyah",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
	"Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
	"Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
	"Phillips", "Evans", "Turner", "Torres", "Parker", "Collins", "Edwards", "Stewart", "Flores", "Morris",
	"Nguyen", "Murphy", "Rivera", "Cook", "Rogers", "Morgan", "Peterson", "Cooper", "Reed", "Bailey",
	"Bell", "Gomez", "Kelly", "Howard", "Ward", "Cox", "Diaz", "Richardson", "Wood", "Watson",
	"Brooks", "Bennett", "Gray", "James", "Reyes", "Cruz", "Hughes", "Price", "Myers", "Long",
	"Foster", "Sanders", "Ross", "Morales", "Powell", "Sullivan", "Russell", "Ortiz", "Jenkins", "Gutierrez",
	"Perry", "Butler", "Barnes", "Fisher", "Henderson", "Coleman", "Simmons", "Patterson", "Jordan", "Reynolds",
}
