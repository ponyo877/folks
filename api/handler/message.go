package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/entity"
	"github.com/ponyo877/folks/usecase/message"
)

// MakeMessageHandlers
func MakeMessageHandlers(e *echo.Echo, service message.UseCase) {
	e.GET("/v1/message", GetMessage(service))
}

// GetMessage
func GetMessage(service message.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Info("WebSocketの接続を開始します")
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Errorf("WebSocketの接続に失敗しました: %v", err)
			return err
		}
		messageChannel := make(chan *entity.Message)
		service.Subscribe(messageChannel)

		session := entity.NewSession(conn)
		// go writeMessage(session, messageChannel)
		go func() {
			for {
				select {
				case message := <-messageChannel:
					if err := session.Conn.WriteMessage(websocket.TextMessage, []byte(message.MessageText())); err != nil {
						log.Errorf("WebSocketのメッセージの送信に失敗しました: %v", err)
					}
				}
			}
		}()
		// go readMessage(session)
		go func() {
			for {
				_, messageText, err := session.Conn.ReadMessage()
				if err != nil {
					log.Errorf("WebSocketのメッセージの受信に失敗しました: %v", err)
				}
				message := entity.Message{
					ID:     entity.UID{},
					UserID: entity.UID{},
					RoomID: entity.UID{},
					Message: entity.Text{
						Value: string(messageText),
					},
					CreatedAt: time.Time{},
				}
				service.Publish(message)
				fmt.Printf("%s\n", message)
			}
		}()
		return nil
	}
}

// writeMessage
func writeMessage(session *entity.Session, messageChannel chan *entity.Message) {
	defer session.Conn.Close()
	for {
		select {
		case message := <-messageChannel:
			if err := session.Conn.WriteMessage(websocket.TextMessage, []byte(message.MessageText())); err != nil {
				log.Errorf("WebSocketのメッセージの送信に失敗しました: %v", err)
			}
		}
	}
}

// readMessage
func readMessage(session *entity.Session) {
	defer session.Conn.Close()
	for {
		_, message, err := session.Conn.ReadMessage()
		if err != nil {
			log.Errorf("WebSocketのメッセージの受信に失敗しました: %v", err)
		}

		fmt.Printf("%s\n", message)
	}
}
