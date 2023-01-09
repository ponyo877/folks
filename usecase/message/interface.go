package message

import "github.com/ponyo877/folks/entity"

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	Publish(binary []byte) error
	Subscribe(f func([]byte)) error
	Append(ID entity.UID, message *entity.Message) error
	ListRecent(ID entity.UID, size int64) ([]*entity.Message, error)
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Publish(message *entity.Message) error
	Subscribe(messageChannel chan *entity.Message) error
	ListRecent(ID entity.UID) ([]*entity.Message, error)
}
