package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	bookcrud "github.com/qazaqpyn/bookCRUD"
	"github.com/qazaqpyn/bookCRUD/pkg/handler"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"github.com/qazaqpyn/bookCRUD/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewMongodb(viper.GetString("mongo.name"))
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(bookcrud.Server)
	//Graceful shutdown implementation
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("BOOKCRUD started\n")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("BOOKCRUD shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while shutting down server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
