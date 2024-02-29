package server

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log-snare/web/data"
	"log-snare/web/service"
	"net/http"
	"strconv"
)

type Server struct {
	Debug bool
}

func Run(configFile string, debug bool, resetDb bool) error {

	if len(configFile) == 0 {
		return errors.New("must provide a config path")
	}

	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"logsnare.log", "stdout"},
		ErrorOutputPaths: []string{"logsnare.log", "stdout"},
		InitialFields: map[string]interface{}{
			"program": "log-snare",
			"version": 0.1,
		},
	}
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("unable to build logger: %v", err)
	}

	// data init
	db, err := gorm.Open(sqlite.Open("logsnare.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if resetDb {
		if err := db.Migrator().DropTable(&data.User{}, &data.Employee{}, &data.Company{}, &data.SettingValue{}); err != nil {
			panic("failed to drop tables")
		}
	}

	err = db.AutoMigrate(&data.User{}, &data.Employee{}, &data.Company{}, &data.SettingValue{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = seedDBIfNeeded(db)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Set up services
	userService := service.NewUserService(db, logger)
	dashboardService := service.NewDashboardService(db, logger)
	employeeService := service.NewEmployeeService(db, logger)
	settingsService := service.NewSettingsService(db, logger)

	// Set up handlers
	userHandler := NewUserHandler(userService, settingsService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	employeeHandler := NewEmployeeHandler(employeeService, settingsService, userService)
	settingsHandler := NewSettingsHandler(settingsService, userService)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("LogSnareSession", store))
	gob.Register(data.UserSafe{})

	// Load the template files
	r.LoadHTMLGlob("../ui/templates/*")

	// service UI assets for anyone
	r.Static("assets/css", "../ui/assets/css")
	r.Static("assets/js", "../ui/assets/js")
	r.Static("assets/img", "../ui/assets/img")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.POST("/login", userHandler.Login)

	authRoutes := r.Group("/app")
	authRoutes.Use(AuthMiddleware())

	// application endpoints begin
	// views
	authRoutes.GET("/dashboard", dashboardHandler.Dashboard)
	authRoutes.GET("/employees/:id", employeeHandler.Employees)
	authRoutes.GET("/settings", settingsHandler.Settings)
	authRoutes.GET("/users/:id", userHandler.Users)
	authRoutes.GET("/impersonate/:id", userHandler.Impersonate)
	authRoutes.GET("/docs", settingsHandler.Docs)

	// api
	authRoutes.GET("/enable-admin", userHandler.EnableAdmin)
	authRoutes.GET("/disable-admin", userHandler.DisableAdmin)

	// educational endpoints that change application behavior.
	authRoutes.GET("/enable-validation", settingsHandler.EnableValidation)
	authRoutes.GET("/disable-validation", settingsHandler.DisableValidation)

	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func seedDBIfNeeded(db *gorm.DB) error {
	var count int64
	db.Model(&data.Company{}).Count(&count)
	if count == 0 {

		logSnareUsers := []data.User{
			data.CreateUserWithPassword("gopher", 2, true),
			data.CreateUserWithPassword("gophmin", 1, false),
		}

		var logSnareEmployees []data.Employee
		for i := 0; i < 20; i++ {
			employee := data.GenerateEmployee("logsnare.local")
			logSnareEmployees = append(logSnareEmployees, employee)
		}

		logSnareCompany := data.Company{
			Name:      "LogSnare",
			Users:     logSnareUsers,
			Employees: logSnareEmployees,
		}
		if err := db.Create(&logSnareCompany).Error; err != nil {
			log.Fatal("failed to create logsnare company:", err)
		}

		acmeUsers := []data.User{
			data.CreateUserWithPassword("acme-admin", 1, false),
			data.CreateUserWithPassword("acme-user", 2, false),
		}

		var acmeEmployees []data.Employee
		for i := 0; i < 13; i++ {
			employee := data.GenerateEmployee("acme.local")
			acmeEmployees = append(acmeEmployees, employee)
		}

		acmeCompany := data.Company{
			Name:      "Acme",
			Users:     acmeUsers,
			Employees: acmeEmployees,
		}
		if err := db.Create(&acmeCompany).Error; err != nil {
			log.Fatal("failed to create acme company:", err)
		}

		// add settings
		db.Create(&data.SettingValue{
			Key:   "1",
			Value: false,
		})
		db.Create(&data.SettingValue{
			Key:   "2",
			Value: false,
		})
		db.Create(&data.SettingValue{
			Key:   "3",
			Value: false,
		})

		log.Println("Database seeded")
	}

	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		if !user.(data.UserSafe).Active {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()

	}
}

func ManageUserCompanyIdentifier(companyId int) string {
	target := "CompanyId:" + strconv.Itoa(companyId)
	return base64.StdEncoding.EncodeToString([]byte(target))
}
