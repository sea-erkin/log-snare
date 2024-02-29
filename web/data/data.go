package data

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func SetupDB(resetDb bool) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("logsnare.db"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("failed to open db file: %w", err)
	}

	if resetDb {
		if err := db.Migrator().DropTable(&User{}, &Employee{}, &Company{}, &SettingValue{}); err != nil {
			return db, fmt.Errorf("failed to drop tables : %w", err)
		}
	}

	err = db.AutoMigrate(&User{}, &Employee{}, &Company{}, &SettingValue{})
	if err != nil {
		return db, fmt.Errorf("failed to auto migrate database: %w", err)
	}
	err = seedDBIfNeeded(db)
	if err != nil {
		return db, fmt.Errorf("failed to seed database: %w", err)
	}

	return db, nil
}

func seedDBIfNeeded(db *gorm.DB) error {
	var count int64
	db.Model(&Company{}).Count(&count)
	if count == 0 {

		logSnareUsers := []User{
			createUserWithPassword("gopher", 2, true),
			createUserWithPassword("gophmin", 1, false),
		}

		var logSnareEmployees []Employee
		for i := 0; i < 20; i++ {
			employee := GenerateEmployee("logsnare.local")
			logSnareEmployees = append(logSnareEmployees, employee)
		}

		logSnareCompany := Company{
			Name:      "LogSnare",
			Users:     logSnareUsers,
			Employees: logSnareEmployees,
		}
		if err := db.Create(&logSnareCompany).Error; err != nil {
			log.Fatal("failed to create logsnare company:", err)
		}

		acmeUsers := []User{
			createUserWithPassword("acme-admin", 1, false),
			createUserWithPassword("acme-user", 2, false),
		}

		var acmeEmployees []Employee
		for i := 0; i < 13; i++ {
			employee := GenerateEmployee("acme.local")
			acmeEmployees = append(acmeEmployees, employee)
		}

		acmeCompany := Company{
			Name:      "Acme",
			Users:     acmeUsers,
			Employees: acmeEmployees,
		}
		if err := db.Create(&acmeCompany).Error; err != nil {
			log.Fatal("failed to create acme company:", err)
		}

		// add settings
		db.Create(&SettingValue{
			Key:   "1",
			Value: false,
		})
		db.Create(&SettingValue{
			Key:   "2",
			Value: false,
		})
		db.Create(&SettingValue{
			Key:   "3",
			Value: false,
		})

		log.Println("Database seeded")
	}

	return nil
}
