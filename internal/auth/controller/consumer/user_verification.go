package consumer

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/consumer/dto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type UserVerificationCallback struct {
	logger  *zap.SugaredLogger
	dbRedis *redis.Client
}

func NewUserVerificationCallback(logger *zap.SugaredLogger, dbRedis *redis.Client) *UserVerificationCallback {
	return &UserVerificationCallback{
		logger:  logger,
		dbRedis: dbRedis,
	}
}

func (c *UserVerificationCallback) Callback(message <-chan *sarama.ConsumerMessage, error <-chan *sarama.ConsumerError) {
	for {
		select {
		case msg := <-message:
			var userCode dto.UserCode

			err := json.Unmarshal(msg.Value, &userCode)
			if err != nil {
				c.logger.Errorf("failed to unmarshall record value err: %v", err)
			} else {
				c.logger.Infof("user code: %s", userCode)

				// save to database
				err = c.dbRedis.Set(context.Background(), userCode.Key, userCode.Code, 2*time.Minute).Err()
				if err != nil {
					c.logger.Errorf("failed to save record value in redis err: %v", err)
				}
			}
		case err := <-error:
			c.logger.Errorf("failed consume err: %v", err)
		}
	}
}
