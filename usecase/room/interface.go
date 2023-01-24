package room

import "github.com/ponyo877/folks/entity"

// Reader interface
type Reader interface {
	ListRecent(roomID entity.UID, size int64) ([]*entity.Message, error)
	ListRoom() ([]*entity.Room, error)
	GetRoom(roomID entity.UID) (*entity.Room, error)
}

// Writer interface
type Writer interface {
	Publish(roomID entity.UID, binary []byte) error
	Subscribe(roomID entity.UID, f func([]byte)) error
	Append(roomID entity.UID, message *entity.Message) error
	CreateRoom(room *entity.Room) error
	AddUser(session *entity.Session) error
	RemoveUser(session *entity.Session) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Publish(roomID entity.UID, message *entity.Message) error
	ConnectRoom(session *entity.Session) (chan *entity.Message, error)
	DisconnectRoom(session *entity.Session) error
	ListRecent(roomID entity.UID) ([]*entity.Message, error)
	ListRoom() ([]*entity.Room, error)
	CreateRoom(room *entity.Room) error
	GetRoom(roomID entity.UID) (*entity.Room, error)
	WriteMessage(session *entity.Session) error
	ReadMessage(session *entity.Session) error
}
