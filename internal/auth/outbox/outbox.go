package outbox

import (
	"context"
	"encoding/json"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/consumer/dto"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
	"github.com/nanmenkaimak/github-gist/internal/auth/repository"
	"github.com/nanmenkaimak/github-gist/internal/kafka"
	"go.uber.org/zap"
	"sync"
	"time"
)

type KafkaOutbox struct {
	dbRepo        repository.Repository
	kafkaProducer *kafka.Producer
	wg            sync.WaitGroup
	logger        *zap.SugaredLogger
	workers       int
	interval      time.Duration
}

func NewKafkaOutbox(dbRepo repository.Repository, kafkaProducer *kafka.Producer,
	logger *zap.SugaredLogger, workers int, interval time.Duration) *KafkaOutbox {
	return &KafkaOutbox{
		dbRepo:        dbRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		workers:       workers,
		interval:      interval,
	}
}

func (k *KafkaOutbox) Run() {
	ticker := time.NewTicker(k.interval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			messages, err := k.dbRepo.GetNotProcessedMessages()
			if err != nil {
				k.logger.Errorf("get not processed messages err: %v", err)
			}

			jobs := make(chan entity.Message, len(messages))

			for i := 1; i <= k.workers; i++ {
				k.wg.Add(1)
				go k.worker(i, jobs)
			}

			for _, event := range messages {
				jobs <- event
			}

			close(jobs)
			k.wg.Wait()

			err = k.dbRepo.DeleteMessage()
			if err != nil {
				k.logger.Errorf("deleting processed messages err: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (k *KafkaOutbox) worker(id int, jobs <-chan entity.Message) {
	defer k.wg.Done()

	for event := range jobs {
		msg := dto.UserCode{
			Code: event.Code,
			Key:  event.Key,
		}

		b, err := json.Marshal(&msg)
		if err != nil {
			k.logger.Errorf("worker %d, failed to marshall UserCode err: %v", id, err)
		}

		k.kafkaProducer.ProduceMessage(b)
	}
}
