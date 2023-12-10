package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"grates/docs"
	"grates/internal/handler"
	"grates/internal/repository"
	"grates/internal/service"
	"grates/pkg/app"
	repoConf "grates/pkg/repository"
	"grates/pkg/server"
	"os"
)

// @title Grates API
// @version 1.0
// @description API Server for Grates social network app

// @contact.name   Yaroslav Molodcov
// @contact.email  iam@it-yaroslav.ru

// @basePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	//logrus.SetFormatter(new(logrus.JSONFormatter))

	configFileName := "config"
	envFileName := ".env"

	isProd := flag.Bool("prod", false, "if the project is launched for production")
	flag.Parse()

	if *isProd {
		configFileName = "config.prod"
		envFileName = ".prod.env"
	}

	if err := godotenv.Load(envFileName); err != nil {
		logrus.Fatalf("error loading env vars: %s", err.Error())
	}

	if err := app.InitConfig(configFileName); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))

	// PosgtgreSQL connect
	db, err := repoConf.NewPostgresDB(repoConf.PSQLConfig{
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
	defer func() { db.Close() }()

	// Redis connect
	rdb := repoConf.NewRedisDB(repoConf.RedisConfig{Addr: viper.GetString("rdb.addr")})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logrus.Fatalf("failed to initialize redis db: %s", err.Error())
	}
	logrus.Info("Redis DB Connected")
	defer func() { rdb.Close() }()

	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos, service.Deps{
		SigingKey:       os.Getenv("JWT_SIGING_KEY"),
		PasswordSalt:    os.Getenv("PASSWORD_SALT"),
		AccessTokenTTL:  viper.GetDuration("auth.accessTokenTTL"),
		RefreshTokenTTL: viper.GetDuration("auth.refreshTokenTTL"),
	})
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	logrus.Info("Grates Server Started")
}
