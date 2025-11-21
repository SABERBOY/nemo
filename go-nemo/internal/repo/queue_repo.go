package repo

import (
	"context"
	"netease-kit/nemo/internal/db"
	"time"

	"github.com/go-redis/redis/v8"
)

type QueueRepository struct{}

func NewQueueRepository() *QueueRepository {
	return &QueueRepository{}
}

func (r *QueueRepository) SendDelayMessage(topic string, message string, delay int64) error {
	// delay is in milliseconds
	score := float64(time.Now().UnixMilli() + delay)
	return db.RDB.ZAdd(context.Background(), topic, &redis.Z{
		Score:  score,
		Member: message,
	}).Err()
}

func (r *QueueRepository) CancelDelayMessage(topic string, message string) error {
	return db.RDB.ZRem(context.Background(), topic, message).Err()
}
