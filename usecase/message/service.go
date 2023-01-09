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
func (s *Service) Publish(message *entity.Message) error {
	messageBinary, err := entity.EncodeMessage(message)
	if err != nil {
		return err
	}
	testRoomUID, _ := entity.StringToID("00000000-0000-0000-0000-000000000001")
	if err := s.repository.Append(testRoomUID, message); err != nil {
		return err
	}
	return s.repository.Publish(messageBinary)
}

// Subscribe
func (s *Service) Subscribe(messageChannel chan *entity.Message) error {
	return s.repository.Subscribe(func(binary []byte) {
		message, err := entity.DecodeMessage(binary)
		if err != nil {
			message = &entity.ErrorMessage
		}
		messageChannel <- message
	})
}

func (s *Service) ListRecent(ID entity.UID) ([]*entity.Message, error) {
	return s.repository.ListRecent(ID, 100)
}
