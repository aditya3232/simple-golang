package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"simple-golang/config"
	inboundadapterecho "simple-golang/internal/adapter/inbound/echo"
	outboundadapterpostgres "simple-golang/internal/adapter/outbound/postgres/repository"
	"simple-golang/internal/domain/service"
	"simple-golang/util"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func RunServer() {
	cfg := config.NewConfig()
	redisConfig := cfg.RedisConfig()

	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-1] Failed to connect Postgres: %v", err)
		return
	}

	if cfg.App.AppPort == "" {
		cfg.App.AppPort = os.Getenv("APP_PORT")
	}
	appPort := ":" + cfg.App.AppPort

	userRepo := outboundadapterpostgres.NewUserRepository(db.DB)

	jwtService := service.NewJwtService(cfg)
	userService := service.NewUserService(userRepo, jwtService, redisConfig)

	e := echo.New()
	e.Use(middleware.CORS())
	e.HideBanner = true
	e.Use(middleware.Recover())

	customValidator := util.NewValidator(db.DB)
	if err := en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator); err != nil {
		log.Fatalf("[RunServer-2] %v", err)
		return
	}
	e.Validator = customValidator

	mid := inboundadapterecho.NewMiddlewareAdapter(cfg, redisConfig, jwtService)
	pingHandler := inboundadapterecho.NewPingHandler()
	userHandler := inboundadapterecho.NewUserHandler(userService)

	inboundadapterecho.InitRoutes(e, mid, userHandler, pingHandler)

	go func() {
		log.Infof("[RunServer-3] Server starting at %s", appPort)
		if err := e.Start(appPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[RunServer-6] Server start failed: %v", err)
		}
	}()

	// === GRACEFUL SHUTDOWN ===
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Infof("[RunServer-4] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("[RunServer-5] Server forced to shutdown: %v", err)
	}

	log.Infof("[RunServer-6] Server exited properly")
}
