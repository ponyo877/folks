package handler

import (
	"net/http"
	"time"

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
		go func() {
			messages, err := service.ListRecent(entity.ErrorUID)
			if err != nil {
				log.Errorf("ListRecentに失敗しました: %v", err)
			}
			for _, message := range messages {
				messageChannel <- message
			}
		}()

		session := entity.NewSession(conn)
		isClosed := false
		isDone := make(chan struct{}, 1)

		session.Conn.SetCloseHandler(func(code int, text string) error {
			isClosed = true
			isDone <- struct{}{}
			return nil
		})

		// サーバ->クライアント
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer func() {
				ticker.Stop()
				session.Conn.Close()
			}()
			for {
				select {
				case message := <-messageChannel:
					log.Infof("message: %v", message)
					messageResponcePresenter := presenter.MarshalMessage(message)
					if err != nil {
						log.Errorf("PickMessageに失敗しました: %v", err)
						return
					}
					if err := session.Conn.WriteJSON(&messageResponcePresenter); err != nil {
						log.Errorf("WebSocketのメッセージの書込に失敗しました: %v", err)
						return
					}
				case <-isDone:
					return
				case <-ticker.C:
					if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						log.Errorf("Pingに失敗しました :%v", err)
						return
					}
				}
			}
		}()
		// クライアント->サーバ
		go func() {
			defer session.Conn.Close()
			for {
				if isClosed {
					return
				}
				var messageRequestPresenter presenter.MessageRequestPresenter
				if err := session.Conn.ReadJSON(&messageRequestPresenter); err != nil {
					log.Infof("WebSocketのメッセージの受信に失敗しました: %v", err)
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Errorf("WebSocketのメッセージの受信が想定外の原因で失敗しました: %v", err)
					}
					return
				}
				message := presenter.UnmarshalMessage(&messageRequestPresenter)
				if err := service.Publish(message); err != nil {
					log.Errorf("WebSocketのメッセージの出版に失敗しました: %v", err)
				}
			}
		}()
		return nil
	}
}
