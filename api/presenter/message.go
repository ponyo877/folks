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
	ID        string `json:"id"`
	UserName  string `json:"userName"`
	UserHash  string `json:"userHash"`
	RoomName  string `json:"roomName"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"createdAt"`
}

// UnmarshalMessage
func UnmarshalMessage(messageRequestPresenter *MessageRequestPresenter) *entity.Message {
	message := &entity.Message{
		ID:        entity.NewUID(),
		UserName:  entity.StringToText(messageRequestPresenter.UserName),
		UserHash:  entity.NewHash(messageRequestPresenter.IPAddress),
		RoomName:  entity.StringToText(messageRequestPresenter.RoomName),
		Message:   entity.StringToText(messageRequestPresenter.Message),
		CreatedAt: time.Now(),
	}
	return message
}

// MarshalMessage
func MarshalMessage(message *entity.Message) MessageResponcePresenter {
	messagePresenter := MessageResponcePresenter{
		ID:        message.ID.String(),
		UserName:  message.UserName.String(),
		UserHash:  message.UserHash.String(),
		RoomName:  message.RoomName.String(),
		Message:   message.MessageText(),
		CreatedAt: message.CreatedAt.UnixMilli(),
	}
	return messagePresenter
}
