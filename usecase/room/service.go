package room

import (
	"github.com/ponyo877/folks/entity"
)

// Service Room usecase
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
func (s *Service) Publish(roomID entity.UID, message *entity.Message) error {
	messageBinary, err := entity.EncodeMessage(message)
	if err != nil {
		return err
	}
	if err := s.repository.Append(roomID, message); err != nil {
		return err
	}
	return s.repository.Publish(roomID, messageBinary)
}

// Subscribe
func (s *Service) Subscribe(roomID entity.UID, messageChannel chan *entity.Message) error {
	return s.repository.Subscribe(
		roomID,
		func(binary []byte) {
			message, err := entity.DecodeMessage(binary)
			if err != nil {
				message = &entity.ErrorMessage
			}
			messageChannel <- message
		},
	)
}

// ListRecent
func (s *Service) ListRecent(roomID entity.UID) ([]*entity.Message, error) {
	return s.repository.ListRecent(roomID, 100)
}

// ListRoom
func (s *Service) ListRoom() error {
	return nil
}

// CreateRoom
func (s *Service) CreateRoom(displayName entity.DisplayName) error {
	return nil
}
