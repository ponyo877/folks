package entity

import (
	"github.com/gorilla/websocket"
)

// Session
type Session struct {
	UserDisplayName DisplayName
	RoomID          UID
	Conn            *websocket.Conn // 本来はフレームワークをentityに持ち込むべきじゃない
	IsDone          chan struct{}
	IsClosed        bool
}

func NewSession(displayName DisplayName, roomID UID, conn *websocket.Conn) *Session {
	session := &Session{
		UserDisplayName: displayName,
		RoomID:          roomID,
		Conn:            conn,
	}
	return session
}

func (s *Session) Close() {
	s.IsClosed = true
	s.IsDone <- struct{}{}
	s.Conn.Close()
}
