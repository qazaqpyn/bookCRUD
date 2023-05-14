package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	bookcrud "github.com/qazaqpyn/bookCRUD"
	"github.com/qazaqpyn/bookCRUD/pkg/handler"
	grpc_client "github.com/qazaqpyn/bookCRUD/pkg/handler/grpc"
	rabbitmq "github.com/qazaqpyn/bookCRUD/pkg/handler/rabbitMQ"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"github.com/qazaqpyn/bookCRUD/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//	@title			Book CRUD API
//	@version		1.0
//	@description	This is a simple server for CRUD operations on book.

//	@contact.name	API Support
//	@contact.email	alimkali.alizhan@gmail.com

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Token-based authentication
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

	auditClient, err := grpc_client.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("RABBITMQ")

	queueServer, err := rabbitmq.NewServer(uri)
	if err != nil {
		log.Fatal(err)
	}

	services := service.NewService(repos, queueServer, auditClient)

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

	if err := auditClient.CloseConnection(); err != nil {
		logrus.Errorf("error occured while shutting down gRPC client: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
