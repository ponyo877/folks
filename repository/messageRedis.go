package repository

import (
	"context"

	"github.com/ponyo877/folks/entity"
)

// Append
func (r *MessageRepository) Append(roomID entity.UID, message *entity.Message) error {
	encodedMessage, err := entity.EncodeMessage(message)
	if err != nil {
		return err
	}
	return r.kvs.RPush(context.Background(), roomID.String(), encodedMessage).Err()
}

// ListRecent
func (r *MessageRepository) ListRecent(roomID entity.UID, size int64) ([]*entity.Message, error) {
	preDecodeMessages, err := r.kvs.LRange(context.Background(), roomID.String(), -size, -1).Result()
	if err != nil {
		return nil, nil
	}
	var messageList []*entity.Message
	for _, preDecodeMessage := range preDecodeMessages {
		message, err := entity.DecodeMessage([]byte(preDecodeMessage))
		if err != nil {
			return nil, nil
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}
