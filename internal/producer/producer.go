package producer

import (
	"context"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"strconv"
	"time"
)

const (
	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second
	writerRequiredAcks = -1
	writerMaxAttempts  = 3
	writerBatchSize    = 1
	writerBatchTimeout = 10 * time.Millisecond
)

type Producer struct {
	log    logger.Logger
	cfg    *config.Config
	Writer *kafka.Writer
}

// NewProducer constructor
func NewProducer(cfg *config.Config, log logger.Logger) *Producer {
	return &Producer{
		log: log,
		cfg: cfg,
	}
}

// GetNewKafkaWriter Create new kafka writer
func (p *Producer) GetNewKafkaWriter(topic string) *kafka.Writer {
	// initialize the writer with the broker addresses, and the topic
	w := &kafka.Writer{
		Addr:  kafka.TCP(p.cfg.Brokers...),
		Topic: topic,
		// LeastBytes is a Balancer implementation that routes messages to the partition that has received the least amount of data.
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		BatchSize:    writerBatchSize,
		// no matter what happens, write all pending messages
		// every 10 milliseconds
		BatchTimeout: writerBatchTimeout,
		Logger:       kafka.LoggerFunc(p.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(p.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

func (p *Producer) WriteMessages(duration int, messagesCount int) error {
	ctx := context.Background()
	// initialize a counter
	i := 0
	for k := 0; k < duration; k++ {

		// start loop with timeout a second
		// only number in messagesCount of operations per second can be reached
		for stay, counter, start, timeout := true, 0, time.Now(), time.After(time.Second); stay; {
			select {
			case <-timeout:
				p.log.Infof("%v messages were sent per second", counter)
				stay = false
				continue
			default:
				// each kafka message has a key and value. The key is used
				// to decide which partition (and consequently, which broker)
				// the message gets published on
				err := p.Writer.WriteMessages(ctx, kafka.Message{
					Key: []byte(strconv.Itoa(i)),
					// create an arbitrary message payload for the value
					Value: []byte("this is message-" + strconv.Itoa(i)),
				})
				if err != nil {
					p.log.Errorf("could not write message: %e - %T", err, err)
					return err
				}
				i++

				// increment local counter
				counter++
				if counter == messagesCount {
					currentTime := time.Now()
					pastTime := currentTime.Sub(start)
					result := time.Second - pastTime
					p.log.Info("remaining time in second - ", result)
					if result > 0 {
						time.Sleep(result)
						p.log.Infof("%v were sent per second", counter)
						stay = false
						continue
					} else {
						p.log.Infof("%v were sent per second", counter)
						stay = false
						continue
					}
				}
			}
		}
	}
	return nil
}
