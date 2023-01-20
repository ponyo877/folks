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

// ConnectRoom
func (s *Service) ConnectRoom(roomID entity.UID) (chan *entity.Message, error) {
	if _, err := s.repository.GetRoom(roomID); err != nil {
		return nil, err
	}
	messageChannel := make(chan *entity.Message)
	s.repository.Subscribe(
		roomID,
		func(binary []byte) {
			message, err := entity.DecodeMessage(binary)
			if err != nil {
				message = &entity.ErrorMessage
			}
			messageChannel <- message
		},
	)
	return messageChannel, nil
}

// ListRecent
func (s *Service) ListRecent(roomID entity.UID) ([]*entity.Message, error) {
	return s.repository.ListRecent(roomID, 100)
}

// ListRoom
func (s *Service) ListRoom() ([]*entity.Room, error) {
	return s.repository.ListRoom()
}

// CreateRoom
func (s *Service) CreateRoom(room *entity.Room) error {
	return s.repository.CreateRoom(room)
}

// GetRoom
func (s *Service) GetRoom(roomID entity.UID) (*entity.Room, error) {
	return s.repository.GetRoom(roomID)
}
