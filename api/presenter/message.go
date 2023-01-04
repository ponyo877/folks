package presenter

import (
	"time"

	"github.com/ponyo877/folks/entity"
)

// MessageRequestPresenter
type MessageRequestPresenter struct {
	UserName  string `json:"userName"`
	IPAddress string `json:"ipAddress"`
	RoomName  string `json:"roomName"`
	Message   string `json:"message"`
}

// MessageResponcePresenter
type MessageResponcePresenter struct {
	ID        string    `json:"id"`
	UserName  string    `json:"userName"`
	UserHash  string    `json:"userHash"`
	RoomName  string    `json:"roomName"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

// UnmarshalMessage
func UnmarshalMessage(messageRequestPresenter MessageRequestPresenter, message *entity.Message) error {
	message = &entity.Message{
		ID:        entity.NewUID(),
		UserName:  entity.TextToString(messageRequestPresenter.UserName),
		UserHash:  entity.NewHash(messageRequestPresenter.IPAddress),
		RoomName:  entity.TextToString(messageRequestPresenter.RoomName),
		Message:   entity.TextToString(messageRequestPresenter.Message),
		CreatedAt: time.Now(),
	}
	return nil
}

// MarshalMessage
func MarshalMessage(message *entity.Message) (MessageResponcePresenter, error) {
	messagePresenter := MessageResponcePresenter{
		ID:        message.ID.String(),
		UserName:  message.UserName.String(),
		UserHash:  message.UserHash.String(),
		RoomName:  message.RoomName.String(),
		Message:   message.MessageText(),
		CreatedAt: message.CreatedAt,
	}
	return messagePresenter, nil
}
