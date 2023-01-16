package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/api/presenter"
	"github.com/ponyo877/folks/entity"
	"github.com/ponyo877/folks/usecase/room"
)

// MakeRoomHandlers
func MakeRoomHandlers(e *echo.Echo, service room.UseCase) {
	e.GET("/v1/room/:roomID", ConnectRoom(service))
	e.GET("/v1/room", ListRoom(service))
	e.POST("/v1/room", CreateRoom(service))
}

// ConnectRoom
func ConnectRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomIDStr := c.Param("roomID")
		roomID, err := entity.StringToID(roomIDStr)
		if err != nil {
			roomID, _ = entity.StringToID("12345678-0000-0000-0000-000000000002")
		}
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
		service.Subscribe(roomID, messageChannel)
		// 過去ログの表示
		go func() {
			messages, err := service.ListRecent(roomID)
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
					log.Infof("[サーバ->クライアント]message: %v", message)
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
				log.Infof("[クライアント->サーバ]messageRequestPresenter: %v", messageRequestPresenter)
				message := presenter.UnmarshalMessage(&messageRequestPresenter)
				log.Infof("[クライアント->サーバ]message: %v", message)
				if err := service.Publish(roomID, message); err != nil {
					log.Errorf("WebSocketのメッセージの出版に失敗しました: %v", err)
				}
			}
		}()
		return nil
	}
}

// CreateRoom
func CreateRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomPresenter := new(presenter.Room)
		if err := c.Bind(roomPresenter); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(roomPresenter); err != nil {
			return err
		}
		if err := service.CreateRoom(entity.NewDisplayName(roomPresenter.DisplayName)); err != nil {
			return err
		}
		return nil
	}
}

// ListRoom
func ListRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		return service.ListRoom()
	}
}
