package consumer

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/malkev1ch/kafka-task/internal/repository"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

const (
	batchPostgresSize = 10
	batchTimeout      = 100 * time.Millisecond
)

func (cg *ConsumerGroup) createWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	r *kafka.Reader,
	wg *sync.WaitGroup,
	workerID int,
	repo *repository.PostgresRepository,
) {
	defer wg.Done()
	defer cancel()

	batchPostgres := &pgx.Batch{}
	insertQuery := "INSERT INTO messages(message_offset, topic, partition, consumer, key, value)" +
		" VALUES ($1, $2, $3, $4, $5, $6)"
	batchCounter := 0
	ticker := time.After(batchTimeout)
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			cg.log.Errorf("FetchMessage", err)
			return
		}

		cg.log.Infof(
			"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
			workerID,
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Key),
			string(msg.Value),
		)

		batchPostgres.Queue(insertQuery, msg.Offset, msg.Topic, msg.Partition, workerID, string(msg.Key), string(msg.Value))
		batchCounter++
		select {
		case <-ticker:
			cg.log.Infof("WORKER: %v, insert %v rows in db because of timeout %v\n", workerID, batchCounter, batchTimeout)
			if err := repo.SendBatch(ctx, batchPostgres); err != nil {
				cg.log.Fatalf("error occurred while inserting message data in postgres: %e", err)
			}
			batchCounter = 0
			batchPostgres = &pgx.Batch{}
			ticker = time.After(batchTimeout)
		default:
			if batchCounter >= batchPostgresSize {
				cg.log.Infof("WORKER: %v, insert %v rows in db", workerID, batchCounter)
				if err := repo.SendBatch(ctx, batchPostgres); err != nil {
					cg.log.Fatalf("error occurred while inserting message data in postgres: %e", err)
				}
				batchCounter = 0
				batchPostgres = &pgx.Batch{}
			}
		}
	}
}
