package main

import (
	"github.com/caarlos0/env"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/consumer"
	"github.com/malkev1ch/kafka-task/internal/producer"
	"github.com/malkev1ch/kafka-task/internal/repository"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const (
	messagesCountPerSecond = 2000
	sendTimeSecond         = 3
	numConsumer            = 1
)

func main() {
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("can't parse configs")
	}

	prod, err := producer.NewProducer(cfg)
	if err != nil {
		log.Fatalf("error while creation producer  - %e", err)
	}

	defer func() {
		if err = prod.Writer.Close(); err != nil {
			log.Fatalf("error while closing producer connection - %e", err)
		}
	}()

	log.Info("producer successfully created")
	err = prod.WriteMessages(sendTimeSecond, messagesCountPerSecond)
	if err != nil {
		log.Fatalf("error while sending messages - %e", err)
	}

	repo := repository.NewPostgresRepository(cfg)
	defer repo.DB.Close()

	cons, err := consumer.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("error while creation producer  - %e", err)
	}

	defer func() {
		if err = cons.Reader.Close(); err != nil {
			log.Fatalf("error while closing consumer connection - %e", err)
		}
	}()

	for i := 0; i < numConsumer; i++ {
		go func() {
			err := cons.ReadMessages(repo)
			if err != nil {
				log.Fatalf("can't read messages - %e", err)
				return
			}
		}()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
