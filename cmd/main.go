package main

import (
	"log"

	bookcrud "github.com/qazaqpyn/bookCRUD"
	"github.com/qazaqpyn/bookCRUD/pkg/handler"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"github.com/qazaqpyn/bookCRUD/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(bookcrud.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
