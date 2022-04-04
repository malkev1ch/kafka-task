package producer

import (
	"context"
	"fmt"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	// initialize the writer with the broker addresses, and the topic
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaHost + ":" + cfg.KafkaPort),
		Topic:        cfg.KafkaTopic,
		BatchSize:    1 << 10, // 1KB
		BatchTimeout: 100 * time.Millisecond,
		RequiredAcks: 1,
	}
	return &Producer{Writer: w}, nil
}

func (p *Producer) WriteMessages(duration int, messagesCount int) error {
	ctx := context.Background()
	// initialize a counter
	i := 0
	for k := 0; k < duration; k++ {
		// start loop with timeout a second
		// only number in messagesCount of operations per second can be reached
		for counter, start, timeout := 0, time.Now(), time.After(time.Second); ; {
			select {
			case <-timeout:
				logrus.Info("%v were sent per second", counter)
				break
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
					logrus.Error("could not write message " + err.Error())
					return err
				}

				// log a confirmation once the message is written
				fmt.Println("writes:", i)
				i++

				counter++
				if counter == messagesCount {
					currentTime := time.Now()
					pastTime := currentTime.Sub(start)
					result := pastTime - time.Second
					if result > 0 {
						time.Sleep(result)
						break
					}
				}
			}
		}
	}
	return nil
}
