package consumer

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/repository"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Consumer struct {
	Reader            *kafka.Reader
	batchPostgresSize int
}

const (
	minBytes = 1
	maxBytes = 1e6
)

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaHost + ":" + cfg.KafkaPort},
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.KafkaGroupID,
		MinBytes: minBytes,
		// the kafka library requires you to set the MaxBytes
		// in case the MinBytes are set
		MaxBytes: maxBytes,
		// wait for at most 5 seconds before receiving new data
		MaxWait:     5 * time.Second,
		StartOffset: kafka.LastOffset,
	})

	return &Consumer{
		Reader: r,
	}, nil
}

func (c *Consumer) ReadMessages(repo *repository.PostgresRepository) error {
	batchPostgres := &pgx.Batch{}
	ctx := context.Background()
	insertQuery := "INSERT INTO messages(message_offset, topic, partition, key, value) VALUES ($1, $2, $3, $4, $5)"
	batchCounter := 0
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		logrus.Infof("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		batchPostgres.Queue(insertQuery, msg.Offset, msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
		batchCounter++
		if batchCounter >= c.batchPostgresSize {
			batchCounter = 0
			if err := repo.SendBatch(context.Background(), batchPostgres); err != nil {
				logrus.Fatal("error occurred while inserting message data in postgres")
			}
		}

	}
	return nil
}
