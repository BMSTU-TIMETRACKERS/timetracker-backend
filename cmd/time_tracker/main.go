package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	configTimeTracker "github.com/BMSTU-TIMETRACKERS/timetracker-backend/config/time_tracker"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/config/time_tracker/flags"
	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
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
	/*ctx*/ _ = context.Background()

	e := echo.New()
	services, err := tt.Init(e)

	logger := services.Logger

	if err != nil {
		logger.Error("can not init services")
		return fmt.Errorf("can not init services: %v", err)
	}

	// postgresClient, err := tt.PostgresClient.Init()

	// if err != nil {
	// 	logger.Error("can not connect to Postgres client: %w", err)
	// 	return err
	// } else {
	// 	logger.Info("Success conect to postgres")
	// }

	// redisSessionClient, err := tt.RedisSessionClient.Init()

	// if err != nil {
	// 	logger.Error("can not connect to Redis session client: %w", err)
	// 	return err
	// } else {
	// 	logger.Info("Success conect to redis")
	// }

	httpServer := tt.Server.Init(e)
	server := configTimeTracker.Server{HttpServer: httpServer}
	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
