package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	configTimeTracker "github.com/BMSTU-TIMETRACKERS/timetracker-backend/config/time_tracker"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/config/time_tracker/flags"
	_ "github.com/BMSTU-TIMETRACKERS/timetracker-backend/docs"
	entryDelivery "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/delivery"
	entryRepo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/repository"
	entryUC "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/usecase"
	goalDelivery "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/goal/delivery"
	goalRepo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/goal/repository"
	goalUC "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/goal/usecase"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/middleware"
	projectDelivery "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/delivery"
	projectRepo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/repository"
	projectUC "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/usecase"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config.toml", "path to config file")
}

type TimeTracker struct {
	configTimeTracker.Base
	PostgresClient            flags.PostgresFlags `toml:"postgres-client"`
	RedisSessionClient        flags.RedisFlags    `toml:"redis-client"`
	RedisProjectStorageClient flags.RedisFlags    `toml:"redis-project-storage-client"`
	Server                    flags.ServerFlags   `toml:"server"`
}

func main() {
	timeTracker := TimeTracker{}
	_, err := toml.DecodeFile(configPath, &timeTracker)

	if err != nil {
		log.Fatal(err)
	}

	err = timeTracker.Run()

	if err != nil {
		log.Fatal(err)
		return
	}
}

func (tt TimeTracker) Run() error {
	ctx := context.Background()

	e := echo.New()
	services, err := tt.Init(e)

	logger := services.Logger

	if err != nil {
		logger.Error("can not init services")
		return fmt.Errorf("can not init services: %v", err)
	}

	postgresClient, err := tt.PostgresClient.Init(ctx)
	if err != nil {
		logger.Error("can not connect to Postgres client: %w", err)
		return err
	} else {
		logger.Info("Success connect to postgres")
	}

	// redisSessionClient, err := tt.RedisSessionClient.Init()

	// if err != nil {
	// 	logger.Error("can not connect to Redis session client: %w", err)
	// 	return err
	// } else {
	// 	logger.Info("Success conect to redis")
	// }

	// Репозитории.
	entryRepository := entryRepo.NewRepository(postgresClient)
	projectRepository := projectRepo.NewRepository(postgresClient)
	goalRepository := goalRepo.NewRepository(postgresClient)

	// Usecases.
	entryUsecase := entryUC.NewUsecase(entryRepository)
	projectUsecase := projectUC.NewUsecase(projectRepository, entryRepository)
	goalUsecase := goalUC.NewUsecase(goalRepository)

	// Мидлвары.
	authMW := middleware.NewAuthMiddleware()

	// Регистрация мидлвар.
	e.Use(authMW.Auth)

	// Регистрация обработчиков.
	entryDelivery.RegisterHandlers(e, entryUsecase, logger)
	projectDelivery.RegisterHandlers(e, projectUsecase, logger)
	goalDelivery.RegisterHandlers(e, goalUsecase, logger)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	httpServer := tt.Server.Init(e)
	server := configTimeTracker.Server{HttpServer: httpServer}
	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
