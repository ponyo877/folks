package message

import "github.com/ponyo877/folks/entity"

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	Publish(roomID entity.UID, binary []byte) error
	Subscribe(roomID entity.UID, f func([]byte)) error
	Append(roomID entity.UID, message *entity.Message) error
	ListRecent(roomID entity.UID, size int64) ([]*entity.Message, error)
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Publish(roomID entity.UID, message *entity.Message) error
	Subscribe(roomID entity.UID, messageChannel chan *entity.Message) error
	ListRecent(roomID entity.UID) ([]*entity.Message, error)
}
