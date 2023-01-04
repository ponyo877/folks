package message

import (
	"github.com/ponyo877/folks/entity"
)

// Service Message usecase
type Service struct {
	repository Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

// Publish
func (s *Service) Publish(message entity.Message) error {
	messageBinary, err := message.EncodeMessage()
	if err != nil {
		return err
	}
	s.repository.Publish(messageBinary)
	return nil
}

// Subscribe
func (s *Service) Subscribe(messageChannel chan *entity.Message) error {
	s.repository.Subscribe(func(binary []byte) {
		message, err := entity.DecodeMessage(binary)
		if err != nil {
			message = &entity.ErrorMessage
		}
		messageChannel <- message
	})
	return nil
}
