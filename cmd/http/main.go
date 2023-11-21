package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"grates/internal/handler"
	"grates/internal/repository"
	"grates/internal/service"
	"grates/pkg/server"
	"os"
)

// @title Grates API
// @version 1.0
// @description API Server for Grates social network app

// @contact.name   Yaroslav Molodcov
// @contact.email  iam@it-yaroslav.ru

// @host localhost:8000
// @basePath /

func main() {
	//logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env vars: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// PosgtgreSQL connect
	db, err := repository.NewPostgresDB(repository.PSQLConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	logrus.Info("Postgres DB Connected")

	// Redis connect
	rdb := repository.NewRedisDB(repository.RedisConfig{Addr: viper.GetString("addr")})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logrus.Fatalf("failed to initialize redis db: %s", err.Error())
	}
	logrus.Info("Redis DB Connected")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	logrus.Info("Grates Server Started")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
