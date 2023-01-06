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

type SimpleMessage struct {
	IDStr       string
	UserNameStr string
	UserHashStr string
	RoomNameStr string
	MessageStr  string
	CreatedAt   time.Time
}

// DecodeMessage
func DecodeMessage(binary []byte) (*Message, error) {
	var simpleMessage SimpleMessage
	err := msgpack.Unmarshal(binary, &simpleMessage)
	if err != nil {
		return &Message{}, nil
	}
	ID, err := StringToID(simpleMessage.IDStr)
	if err != nil {
		return &Message{}, nil
	}
	message := &Message{
		ID:        ID,
		UserName:  StringToText(simpleMessage.UserNameStr),
		UserHash:  StringToHash(simpleMessage.UserHashStr),
		RoomName:  StringToText(simpleMessage.RoomNameStr),
		Message:   StringToText(simpleMessage.MessageStr),
		CreatedAt: simpleMessage.CreatedAt,
	}
	return message, err
}

// EncodeMessage
func EncodeMessage(message *Message) ([]byte, error) {
	simpleMessage := &SimpleMessage{
		IDStr:       message.ID.String(),
		UserNameStr: message.UserName.String(),
		UserHashStr: message.UserHash.String(),
		RoomNameStr: message.RoomName.String(),
		MessageStr:  message.Message.String(),
		CreatedAt:   message.CreatedAt,
	}
	return msgpack.Marshal(simpleMessage)
}

// MessageText
func (message *Message) MessageText() string {
	return message.Message.String()
}

var ErrorMessage = Message{
	ID:        ErrorUID,
	UserHash:  NewHash("ErrorMessage"),
	RoomName:  StringToText("ErrorMessage"),
	Message:   StringToText("ErrorMessage"),
	CreatedAt: time.Time{},
}
