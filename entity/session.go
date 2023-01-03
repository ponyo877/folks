package entity

import (
	"github.com/gorilla/websocket"
)

// Session
type Session struct {
	ID     UID
	UserID UID
	RoomID UID
	Conn   *websocket.Conn
}

func NewSession(conn *websocket.Conn) *Session {
	session := new(Session)
	session.Conn = conn
	return session
}
