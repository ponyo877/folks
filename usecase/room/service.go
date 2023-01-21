package room

import (
	"time"

	"github.com/gorilla/websocket"
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

// WriteMessage UIに依存しているためusecase層には移せなさそう
func (s *Service) WriteMessage(session *entity.Session) error {
	// presenterがUIに依存するのでここには書けない
	// ↑が解決すればsession.Connのメソッドをrepositoryに書きたい
	// そうすればapi.handlerのwriteMessageをここに移せる
	// repositoryでwebsocketに送るオブジェクトをmessageをバイナリかしたものにすれば行けそう
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		session.Close()
	}()
	messageChannel, err := s.ConnectRoom(session.RoomID)
	if err != nil {
		return nil
	}
	for {
		select {
		case <-session.IsDone:
			return nil
		case message := <-messageChannel:
			// TODO: WriteMessageをrepository層に新設してmessageをバイナリ化したものを入れる
			if err := session.Conn.WriteMessage(websocket.TextMessage, []byte(message.ID.String())); err != nil {
				return err
			}
		case <-ticker.C:
			if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}

// ReadMessage UIに依存しているためusecase層には移せなさそう
func (s *Service) ReadMessage(session *entity.Session) error {
	// presenterがUIに依存するのでここには書けない
	// ↑が解決すればsession.Connのメソッドをrepositoryに書きたい
	// そうすればapi.handlerのreadMessageをここに移せる
	// repositoryでwebsocketに送るオブジェクトをmessageをバイナリかしたものにすれば行けそう
	defer session.Close()
	for {
		if session.IsClosed {
			return nil
		}
		// TODO: WriteMessageをrepository層に新設してmessageをバイナリ化したものを入れる
		_, _, err := session.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			return nil
		}
		// message.UserName = entity.StringToText(session.UserDisplayName.String())
		// TODO:  &entity.Message{}をmessageに修正
		if err := s.Publish(session.RoomID, &entity.Message{}); err != nil {
		}
	}
}
