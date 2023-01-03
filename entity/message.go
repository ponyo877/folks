package entity

import (
	"time"

	"github.com/vmihailenco/msgpack"
)

type Message struct {
	ID        UID
	UserID    UID
	RoomID    UID
	Message   Text
	CreatedAt time.Time
}

func DecodeMessage(binary []byte) (*Message, error) {
	var message Message
	err := msgpack.Unmarshal(binary, &message)
	return &message, err
}

func (message *Message) EncodeMessage() ([]byte, error) {
	return msgpack.Marshal(message)
}

var ErrorMessage = Message{
	ID:     ErrorUID,
	UserID: ErrorUID,
	RoomID: ErrorUID,
	Message: Text{
		Value: "ErrorMessage",
	},
	CreatedAt: time.Time{},
}

func (message *Message) MessageText() string {
	return message.Message.Value
}
