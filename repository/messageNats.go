package repository

import (
	"github.com/nats-io/nats.go"
	"github.com/ponyo877/folks/entity"
)

// Publish
func (r *MessageRepository) Publish(roomID entity.UID, binary []byte) error {
	return r.mq.Publish(roomID.String(), binary)
}

// Subscribe
func (r *MessageRepository) Subscribe(roomID entity.UID, f func([]byte)) error {
	_, err := r.mq.Subscribe(
		roomID.String(),
		func(m *nats.Msg) {
			f(m.Data)
		},
	)
	return err
}
