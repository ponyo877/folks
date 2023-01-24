package entity

import (
	"github.com/gorilla/websocket"
)

// Session
type Session struct {
	User     User
	RoomID   UID
	Conn     *websocket.Conn // 本来はフレームワークをentityに持ち込むべきじゃない
	IsDone   chan struct{}
	IsClosed bool
}

func NewSession(user User, roomID UID, conn *websocket.Conn) *Session {
	session := &Session{
		User:   user,
		RoomID: roomID,
		Conn:   conn,
	}
	return session
}

func (s *Session) Close() {
	s.IsClosed = true
	s.IsDone <- struct{}{}
	s.Conn.Close()
}
