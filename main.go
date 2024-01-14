package main

import (
	"fmt"
	"github.com/mehmetokdemir/social-media-api/internal/app/user"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"github.com/mehmetokdemir/social-media-api/internal/logger"
	"github.com/mehmetokdemir/social-media-api/internal/postgres"
	"github.com/mehmetokdemir/social-media-api/server"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	appConfig := config.NewConfig()
	appConfig.Print()

	zapLogger, err := logger.NewZapLoggerForEnv(appConfig.AppEnv, 0)
	if err != nil {
		return err
	}

	db, err := postgres.New(appConfig)
	if err != nil {
		return err
	}

	appServer := server.New([]server.Handler{
		prepareUserHandler(db, zapLogger),
	}, appConfig, zapLogger)

	return appServer.Start()
}

func prepareUserHandler(db *gorm.DB, logger *zap.SugaredLogger) *user.HttpHandler {
	userRepository := user.NewRepository(db, logger)
	if err := userRepository.Migration(); err != nil {
		return nil
	}
	userService := user.NewUserService(userRepository, logger)
	userHandler := user.NewHttpHandler(userService, logger)
	return userHandler
}
