package message

import "github.com/ponyo877/folks/entity"

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	Publish(binary []byte) error
	Subscribe(f func([]byte)) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Subscribe(messageChannel chan *entity.Message) error
	Publish(message entity.Message) error
}
