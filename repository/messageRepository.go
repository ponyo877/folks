package repository

import (
	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// MessageRepository repository
type MessageRepository struct {
	mq  *nats.Conn
	kvs *redis.Client
	rdb *gorm.DB
}

// NewMessageRepository create new repository
func NewMessageRepository(mq *nats.Conn, kvs *redis.Client, rdb *gorm.DB) *MessageRepository {
	return &MessageRepository{
		mq,
		kvs,
		rdb,
	}
}
