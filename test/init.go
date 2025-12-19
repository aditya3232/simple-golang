package test

import (
	"log"
	"os"
	"path/filepath"
	"simple-golang/config"
	"simple-golang/internal/domain/service"
	"simple-golang/util"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	inboundadapterecho "simple-golang/internal/adapter/inbound/echo"
	outboundadapterpostgres "simple-golang/internal/adapter/outbound/postgres/repository"
)

func InitViperTestEnv() {
	viper.Reset()

	// CARI PROJECT ROOT
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("[InitViperTestEnv] %v", err)
	}

	// naik sampai ketemu go.mod
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			log.Fatal("[InitViperTestEnv] go.mod not found")
		}
		root = parent
	}

	viper.SetConfigFile(filepath.Join(root, ".env.test"))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("[InitViperTestEnv] %v", err)
	}
}

func InitTestApp() (*echo.Echo, func()) {
	// ===== CONFIG =====
	InitViperTestEnv()
	cfg := config.NewConfig()
	redisConfig := cfg.RedisConfig()

	// ===== DATABASE =====
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[InitTestApp-1] Failed connect postgres: %v", err)
	}

	// ===== REPOSITORY =====
	userRepo := outboundadapterpostgres.NewUserRepository(db.DB)

	// ===== SERVICE =====
	jwtService := service.NewJwtService(cfg)
	userService := service.NewUserService(
		userRepo,
		jwtService,
		redisConfig,
	)

	// ===== ECHO =====
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ===== VALIDATOR =====
	customValidator := util.NewValidator(db.DB)
	if err := en.RegisterDefaultTranslations(
		customValidator.Validator,
		customValidator.Translator,
	); err != nil {
		log.Fatalf("[InitTestApp-2] %v", err)
	}
	e.Validator = customValidator

	// ===== MIDDLEWARE & HANDLER =====
	mid := inboundadapterecho.NewMiddlewareAdapter(cfg, redisConfig, jwtService)
	pingHandler := inboundadapterecho.NewPingHandler()
	userHandler := inboundadapterecho.NewUserHandler(userService)

	// ===== ROUTES =====
	inboundadapterecho.InitRoutes(e, mid, userHandler, pingHandler)

	// ===== CLEANUP =====
	cleanup := func() {
		// Close Postgres
		sqlDB, err := db.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}

		// Close Redis (jika ada client)
		if redisConfig != nil {
			_ = redisConfig.Close()
		}
	}

	return e, cleanup
}
