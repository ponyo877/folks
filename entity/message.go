package entity

import (
	"time"

	"github.com/vmihailenco/msgpack"
)

type Message struct {
	ID        UID
	UserName  Text
	UserHash  Hash
	RoomName  Text
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
	ID:        ErrorUID,
	UserHash:  NewHash("ErrorMessage"),
	RoomName:  TextToString("ErrorMessage"),
	Message:   TextToString("ErrorMessage"),
	CreatedAt: time.Time{},
}

func (message *Message) MessageText() string {
	return message.Message.String()
}
