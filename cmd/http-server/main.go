package main

import (
	"github.com/ikarizxc/http-server/internal/db/postgres"
	"github.com/ikarizxc/http-server/internal/handler"
	"github.com/ikarizxc/http-server/internal/repository"
	"github.com/ikarizxc/http-server/internal/server"
	"github.com/ikarizxc/http-server/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error occured while initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Print("No .env file found")
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     "localhost",
		Port:     "5436",
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("error occured while connecting to database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}
