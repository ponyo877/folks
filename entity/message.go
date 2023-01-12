package entity

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/vmihailenco/msgpack"
)

type Message struct {
	ID        UID
	UserName  Text
	UserHash  Hash
	RoomID    UID
	Message   Text
	CreatedAt time.Time
}

type SimpleMessage struct {
	IDStr       string
	UserNameStr string
	UserHashStr string
	RoomIDStr   string
	MessageStr  string
	CreatedAt   time.Time
}

// DecodeMessage
func DecodeMessage(binary []byte) (*Message, error) {
	var simpleMessage SimpleMessage
	if err := msgpack.Unmarshal(binary, &simpleMessage); err != nil {
		log.Infof("binary is not simpleMessage: %v", err)
		return &Message{}, nil
	}
	ID, err := StringToID(simpleMessage.IDStr)
	if err != nil {
		log.Infof("IDStr is not UUID: %v", err)
		return &Message{}, nil
	}
	roomID, err := StringToID(simpleMessage.RoomIDStr)
	if err != nil {
		log.Infof("RoomIDStr is not UUID: %v", err)
		return &Message{}, nil
	}
	message := &Message{
		ID:        ID,
		UserName:  StringToText(simpleMessage.UserNameStr),
		UserHash:  StringToHash(simpleMessage.UserHashStr),
		RoomID:    roomID,
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
		RoomIDStr:   message.RoomID.String(),
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
	RoomID:    ErrorUID,
	Message:   StringToText("ErrorMessage"),
	CreatedAt: time.Time{},
}
