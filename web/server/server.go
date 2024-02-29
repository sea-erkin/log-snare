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
	"log"
	"log-snare/web/data"
	"log-snare/web/service"
	"net/http"
	"strconv"
)

type Server struct {
	Debug bool
}

func Run(configFile string, debug bool, resetDb bool, listenHost string) error {

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

	db, err := data.SetupDB(resetDb)
	if err != nil {
		logger.Fatal("unable to setup db", zap.Error(err))
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
	// note-unsafe: in prod you'd ideally want to generate a crypto random secure secret for your cookie store.
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
	authRoutes.Use(authMiddleware())

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

	log.Printf("Listening on: http://%s \n", listenHost)

	return r.Run(listenHost)
}

func authMiddleware() gin.HandlerFunc {
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

func manageUserCompanyIdentifier(companyId int) string {
	target := "CompanyId:" + strconv.Itoa(companyId)
	return base64.StdEncoding.EncodeToString([]byte(target))
}
