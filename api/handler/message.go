package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/api/presenter"
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
					log.Infof("MessageResponce: %v", message)
					messageResponcePresenter := presenter.MarshalMessage(message)
					if err != nil {
						log.Errorf("PickMessageに失敗しました: %v", err)
					}
					// log.Infof("MessageResponce: %v", message)
					if err := session.Conn.WriteJSON(&messageResponcePresenter); err != nil {
						log.Errorf("WebSocketのメッセージの送信に失敗しました: %v", err)
					}
					// if err := session.Conn.WriteMessage(websocket.TextMessage, []byte(message.MessageText())); err != nil {
					// 	log.Errorf("WebSocketのメッセージの送信に失敗しました: %v", err)
					// }
				}
			}
		}()
		// go readMessage(session)
		go func() {
			for {
				var messageRequestPresenter presenter.MessageRequestPresenter
				if err := session.Conn.ReadJSON(&messageRequestPresenter); err != nil {
					log.Errorf("WebSocketのメッセージの受信に失敗しました: %v", err)
					continue
				}
				log.Infof("MessageRequestPresenter: %v", messageRequestPresenter)
				message := presenter.UnmarshalMessage(&messageRequestPresenter)
				log.Infof("MessageRequest: %v", message)
				if err := service.Publish(message); err != nil {
					log.Errorf("WebSocketのメッセージの受信に失敗しました: %v", err)
				}
			}
		}()
		return nil
	}
}
