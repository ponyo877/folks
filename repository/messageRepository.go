package repository

import (
	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
)

// MessageRepository repository
type MessageRepository struct {
	mq  *nats.Conn
	kvs *redis.Client
}

// NewMessageRepository create new repository
func NewMessageRepository(mq *nats.Conn, kvs *redis.Client) *MessageRepository {
	return &MessageRepository{
		mq,
		kvs,
	}
}
