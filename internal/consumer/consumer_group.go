package consumer

import (
	"context"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/malkev1ch/kafka-task/internal/repository"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type ConsumerGroup struct {
	Brokers []string
	GroupID string
	log     logger.Logger
	cfg     *config.Config
}

const (
	minBytes               = 10e3 // 10KB
	maxBytes               = 10e6 // 10MB
	queueCapacity          = 100
	heartbeatInterval      = 1 * time.Second
	commitInterval         = 0
	partitionWatchInterval = 5 * time.Second
	maxAttempts            = 3
	dialTimeout            = 3 * time.Minute
)

// NewConsumerGroup constructor
func NewConsumerGroup(
	brokers []string,
	groupID string,
	log logger.Logger,
	cfg *config.Config,
) *ConsumerGroup {
	return &ConsumerGroup{
		Brokers: brokers,
		GroupID: groupID,
		log:     log,
		cfg:     cfg,
	}
}

func (cg *ConsumerGroup) getNewKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(cg.log.Debugf),
		ErrorLogger:            kafka.LoggerFunc(cg.log.Errorf),
		MaxAttempts:            maxAttempts,
		StartOffset:            kafka.LastOffset,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}

func (cg *ConsumerGroup) consumeCreateMessage(ctx context.Context,
	cancel context.CancelFunc,
	groupID string,
	topic string,
	workersNum int,
	repo *repository.PostgresRepository,
) {
	r := cg.getNewKafkaReader(cg.Brokers, topic, groupID)
	defer func() {
		if err := r.Close(); err != nil {
			cg.log.Errorf("r.Close", err)
			cancel()
		}
		cg.log.Infof("consumer group successfully stopped")
	}()

	cg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}
	for i := 0; i < workersNum; i++ {
		wg.Add(1)
		go cg.createWorker(ctx, cancel, r, wg, i, repo)
	}
	wg.Wait()
	cg.log.Infof("consumers stopped")
}

// RunConsumers run kafka consumers
func (cg *ConsumerGroup) RunConsumers(
	ctx context.Context,
	cancel context.CancelFunc,
	messageTopic string,
	workersNum int,
	repo *repository.PostgresRepository) {
	go cg.consumeCreateMessage(ctx, cancel, cg.GroupID, messageTopic, workersNum, repo)
}
