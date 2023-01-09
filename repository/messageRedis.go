package repository

import (
	"context"

	"github.com/ponyo877/folks/entity"
)

// Append
func (r *MessageRepository) Append(ID entity.UID, message *entity.Message) error {
	encodedMessage, err := entity.EncodeMessage(message)
	if err != nil {
		return err
	}
	return r.kvs.RPush(context.Background(), ID.String(), encodedMessage).Err()
}

// ListRecent
func (r *MessageRepository) ListRecent(ID entity.UID, size int64) ([]*entity.Message, error) {
	preDecodeMessages, err := r.kvs.LRange(context.Background(), ID.String(), -size, -1).Result()
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
