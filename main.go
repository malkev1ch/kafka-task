package main

import (
	"context"
	"github.com/caarlos0/env"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/consumer"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/malkev1ch/kafka-task/internal/producer"
	"github.com/malkev1ch/kafka-task/internal/repository"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	messagesCountPerSecond = 100
	sendTimeSecond         = 2
)

func main() {

	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("can't parse configs")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting kafka server")

	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		appLogger.Fatal("cannot connect postgresDB: ", err)
	}
	appLogger.Infof("postgresDB connected")
	defer repo.DB.Close()

	time.Sleep(time.Second * 5)

	consGroup := consumer.NewConsumerGroup(cfg.Brokers, cfg.KafkaGroupID, appLogger, cfg)
	consGroup.RunConsumers(ctx, cancel, repo)

	prod := producer.NewProducer(cfg, appLogger)
	prod.Writer = prod.GetNewKafkaWriter(cfg.KafkaTopic)

	defer func() {
		appLogger.Info("closing producer connection")
		if err := prod.Writer.Close(); err != nil {
			appLogger.Fatalf("error while closing producer connection - %e", err)
		}
	}()

	go func() {
		appLogger.Info("producer starts push message in topic")
		if err := prod.WriteMessages(sendTimeSecond, messagesCountPerSecond); err != nil {
			appLogger.Fatalf("error while sending messages - %e", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
