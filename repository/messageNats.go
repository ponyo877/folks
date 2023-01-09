package repository

import (
	"github.com/nats-io/nats.go"
)

// Publish
func (r *MessageRepository) Publish(binary []byte) error {
	return r.mq.Publish("TestRoom", binary)
}

// Subscribe
func (r *MessageRepository) Subscribe(f func([]byte)) error {
	_, err := r.mq.Subscribe("TestRoom",
		func(m *nats.Msg) {
			f(m.Data)
		},
	)
	return err
}
