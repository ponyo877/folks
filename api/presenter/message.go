package presenter

import (
	"time"

	"github.com/ponyo877/folks/entity"
)

// MessageRequestPresenter
type MessageRequestPresenter struct {
	UserName  string `json:"userName"`
	IPAddress string `json:"ipAddress"`
	RoomID    string `json:"roomID"`
	Message   string `json:"message"`
}

// MessageResponcePresenter
type MessageResponcePresenter struct {
	ID        string `json:"id"`
	UserName  string `json:"userName"`
	UserHash  string `json:"userHash"`
	RoomID    string `json:"roomID"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"createdAt"`
}

// UnmarshalMessage
func UnmarshalMessage(messageRequestPresenter *MessageRequestPresenter) *entity.Message {
	roomID, err := entity.StringToID(messageRequestPresenter.RoomID)
	if err != nil {
		return &entity.Message{}
	}
	message := &entity.Message{
		ID:        entity.NewUID(),
		UserName:  entity.StringToText(messageRequestPresenter.UserName),
		UserHash:  entity.NewHash(messageRequestPresenter.IPAddress),
		RoomID:    roomID,
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
		RoomID:    message.RoomID.String(),
		Message:   message.MessageText(),
		CreatedAt: message.CreatedAt.UnixMilli(),
	}
	return messagePresenter
}

// MarshalMessages
func MarshalMessages(messages []*entity.Message) []MessageResponcePresenter {
	var messageResponcePresenterList []MessageResponcePresenter
	for _, message := range messages {
		messageResponcePresenterList = append(messageResponcePresenterList, MarshalMessage(message))
	}
	return messageResponcePresenterList
}
