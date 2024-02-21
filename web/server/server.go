package server

import (
	"crypto/rand"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log-snare/web/data"
)

type Server struct {
	Debug bool
}

func Run(configFile string, debug bool, resetDb bool) error {

	if len(configFile) == 0 {
		return errors.New("must provide a config path")
	}

	r := gin.Default()

	// data init
	db, err := gorm.Open(sqlite.Open("logsnare.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if resetDb {
		if err := db.Migrator().DropTable(&data.User{}, &data.Company{}); err != nil {
			panic("failed to drop tables")
		}
	}

	err = db.AutoMigrate(&data.User{}, &data.Company{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = seedDBIfNeeded(db)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Load the template files
	r.LoadHTMLGlob("../ui/templates/*")

	// Load UI assets
	r.Static("assets/css", "../ui/assets/css")
	r.Static("assets/js", "../ui/assets/js")
	r.Static("assets/img", "../ui/assets/img")

	// Route for the root path to serve the login page
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})

	r.GET("/users", func(c *gin.Context) {
		c.HTML(200, "users.html", nil)
	})

	// Route to handle POST requests on /login
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		// Here, you would typically authenticate the user.
		// For this example, let's just print the credentials.
		println("Username:", username, "Password:", password)
		// Redirect or respond based on authentication (skipped for simplicity)
		c.JSON(200, gin.H{"status": "submitted"})
	})

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func seedDBIfNeeded(db *gorm.DB) error {
	var count int64
	db.Model(&data.Company{}).Count(&count)
	if count == 0 {

		// create users for LogSnare
		var logSnareUsers []data.User

		logSnareGopher := data.User{
			Username: "gopher",
			Role:     2,
		}
		gopherPass := generatePassword(24)
		logSnareGopher.Password = gopherPass

		logSnareGopherAdmin := data.User{
			Username: "gophmin",
			Role:     1,
		}
		gopherAdminPass := generatePassword(24)
		logSnareGopherAdmin.Password = gopherAdminPass

		logSnareUsers = append(logSnareUsers, logSnareGopher)
		logSnareUsers = append(logSnareUsers, logSnareGopherAdmin)

		log.Println("[DATA] gopher password: ", gopherPass)
		log.Println("[DATA] gophmin password: ", gopherAdminPass)

		// create logsnare company
		company := data.Company{
			Name:  "LogSnare",
			Users: logSnareUsers,
		}

		if err := db.Create(&company).Error; err != nil {
			log.Fatal("failed to create logsnare company")
		}

		// create acme company
		company = data.Company{
			Name: "Acme",
			Users: []data.User{
				{
					Username: "acme-admin",
					Role:     1,
					Password: generatePassword(24),
				},
				{
					Username: "acme-user",
					Role:     2,
					Password: generatePassword(24),
				},
			},
		}

		if err := db.Create(&company).Error; err != nil {
			log.Fatal("failed to create acme company")
		}

		log.Println("Database seeded")
	}

	return nil
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
