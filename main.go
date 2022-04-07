package main

import (
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/consumer"
	"github.com/malkev1ch/kafka-task/internal/handler"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/malkev1ch/kafka-task/internal/producer"
	"github.com/malkev1ch/kafka-task/internal/repository"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("can't parse configs")
	}

	log.Printf("cfg: %+v\n", cfg)

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting kafka server")

	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		appLogger.Fatal("cannot connect postgresDB: ", err)
	}
	appLogger.Info("postgresDB connected")
	defer repo.DB.Close()

	consGroup := consumer.NewConsumerGroup(cfg.Brokers, cfg.KafkaGroupID, appLogger, cfg)
	appLogger.Info("ConsumerGroup created")

	prod := producer.NewProducer(cfg, appLogger)
	appLogger.Info("Producer created")

	defer func() {
		appLogger.Info("closing producer connection")
		if err := prod.Writer.Close(); err != nil {
			appLogger.Fatalf("error while closing producer connection - %e", err)
		}
	}()

	handlers := handler.NewHandler(prod, consGroup, cfg, repo, appLogger)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("messages", handlers.ProduceMessage)
	e.GET("consumer/start", handlers.StartConsumer)
	e.GET("producer/start", handlers.StartProducer)

	// Start server
	go func() {
		e.Logger.Fatal(e.Start(":8080"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
