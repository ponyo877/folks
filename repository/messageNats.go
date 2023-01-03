package repository

import (
	"github.com/nats-io/nats.go"
)

// MessageNats nats repository
type MessageNats struct {
	mq *nats.Conn
}

type MessageNatsPresenter struct {
}

// NewMessageNats create new repository
func NewMessageNats(mq *nats.Conn) *MessageNats {
	return &MessageNats{
		mq: mq,
	}
}

// Publish
func (r *MessageNats) Publish(binary []byte) error {
	return r.mq.Publish("TestRoom", binary)
}

// Subscribe
func (r *MessageNats) Subscribe(f func([]byte)) error {
	_, err := r.mq.Subscribe("TestRoom",
		func(m *nats.Msg) {
			f(m.Data)
		},
	)
	return err
}
