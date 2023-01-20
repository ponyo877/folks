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
	e.POST("/v1/room", CreateRoom(service))
	e.GET("/v1/room", ListRoom(service))
}

// ConnectRoom
func ConnectRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomIDStr := c.Param("roomID")
		roomID, err := entity.StringToID(roomIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		messageChannel, err := service.ConnectRoom(roomID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

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

		conn, err := generateConnection(c.Response(), c.Request())
		if err != nil {
			log.Errorf("WebSocketの接続に失敗しました: %v", err)
		}
		session := entity.NewSession(conn)
		session.Conn.SetCloseHandler(func(code int, text string) error {
			session.IsClosed = true
			session.IsDone <- struct{}{}
			return nil
		})
		go writeToClient(session, messageChannel)
		go writeToServer(session, service, roomID)
		return nil
	}
}

func generateConnection(writer http.ResponseWriter, reader *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return upgrader.Upgrade(writer, reader, nil)
}

func writeToClient(session *entity.Session, messageChannel chan *entity.Message) {
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
			if err := session.Conn.WriteJSON(&messageResponcePresenter); err != nil {
				log.Errorf("WebSocketのメッセージの書込に失敗しました: %v", err)
				return
			}
		case <-session.IsDone:
			return
		case <-ticker.C:
			if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("Pingに失敗しました :%v", err)
				return
			}
		}
	}
}

func writeToServer(session *entity.Session, service room.UseCase, roomID entity.UID) {
	defer session.Conn.Close()
	for {
		if session.IsClosed {
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
}

// CreateRoom
func CreateRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomPresenter := new(presenter.RoomPresenter)
		if err := c.Bind(roomPresenter); err != nil {
			log.Errorf("リクエストの受け取りに失敗しました: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		room := entity.NewRoom(entity.NewDisplayName(roomPresenter.DisplayName))
		if err := service.CreateRoom(room); err != nil {
			log.Errorf("Roomの作成に失敗しました: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, nil)
	}
}

// ListRoom
func ListRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomList, err := service.ListRoom()
		if err != nil {
			log.Errorf("Roomの一覧取得に失敗しました: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, presenter.PickRoomList(roomList))
	}
}
