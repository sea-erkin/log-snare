package server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log-snare/web/data"
	"log-snare/web/service"
	"net/url"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	Debug bool
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func Run(configFile string, debug bool, resetDb bool) error {

	if len(configFile) == 0 {
		return errors.New("must provide a config path")
	}

	ll := lumberjack.Logger{
		Filename:   "logsnare.log",
		MaxSize:    20, //MB
		MaxBackups: 14,
		MaxAge:     1, //days
		Compress:   true,
	}
	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &ll,
		}, nil
	})

	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			TimeKey:    "time",
			EncodeTime: zapcore.EpochMillisTimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{fmt.Sprintf("lumberjack:%s", "logsnare.log"), "stdout"},
		ErrorOutputPaths: []string{fmt.Sprintf("lumberjack:%s", "logsnare.log"), "stdout"},
		InitialFields: map[string]interface{}{
			"program": "log-snare",
			"version": 0.1,
		},
	}
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
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
		if err := db.Migrator().DropTable(&data.User{}, &data.Employee{}, &data.Company{}); err != nil {
			panic("failed to drop tables")
		}
	}

	err = db.AutoMigrate(&data.User{}, &data.Employee{}, &data.Company{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = seedDBIfNeeded(db)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Set up services
	userService := service.NewUserService(db, logger)

	// Set up handlers
	userHandler := NewUserHandler(userService)
	_ = userHandler

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("LogSnareSession", store))
	gob.Register(data.UserSafe{})

	// Load the template files
	r.LoadHTMLGlob("../ui/templates/*")

	// Load UI assets
	r.Static("assets/css", "../ui/assets/css")
	r.Static("assets/js", "../ui/assets/js")
	r.Static("assets/img", "../ui/assets/img")

	r.GET("/users", func(c *gin.Context) {
		c.HTML(200, "users.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})

	r.POST("/login", userHandler.LoginHandler)

	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func seedDBIfNeeded(db *gorm.DB) error {
	var count int64
	db.Model(&data.Company{}).Count(&count)
	if count == 0 {

		logSnareUsers := []data.User{
			data.CreateUserWithPassword("gopher", 2, true),
			data.CreateUserWithPassword("gophmin", 1, true),
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

		log.Println("Database seeded")
	}

	return nil
}
