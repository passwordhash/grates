package main

import (
	"context"
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

const (
	defaultConfigName  = "config"
	defaultEnvFileName = ".env"
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

	envFileName := defaultEnvFileName
	configFileName := defaultConfigName

	envFile := os.Getenv("ENV_FILE")
	if len(envFile) != 0 {
		envFileName = envFile
	}

	if err := godotenv.Load(envFileName); err != nil {
		logrus.Errorf("cannot load env file %s: %s\ntrying lo load .env", envFileName, err.Error())
		envFileName = defaultEnvFileName
		if err := godotenv.Load(); err != nil {
			logrus.Fatalf("error loading env vars: %s", err.Error())
		}
	}

	configFile := os.Getenv("CONFIG_FILE_NAME")
	if len(configFile) != 0 {
		configFileName = configFile
	}

	if err := app.InitConfig(configFileName); err != nil {
		logrus.Errorf("cannot load config file %s: %s\ntrying to load %s", configFileName, err.Error(), defaultConfigName)
		configFileName = defaultConfigName
		if err := app.InitConfig(configFileName); err != nil {
			logrus.Fatalf("error initializing configs: %s", err.Error())
		}
	}

	logrus.Infof("%s env file was loaded", envFileName)
	logrus.Infof("%s config file was loaded", configFileName)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))

	logrus.Info(viper.GetString("host"))

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
		EmailDeps: service.EmailDeps{
			SmtpHost: viper.GetString("email.smtpHost"),
			SmtpPort: viper.GetInt("email.smtpPort"),
			From:     viper.GetString("email.from"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	})
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	logrus.Info("Grates Server Started")
}
