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
	e.GET("/v1/room/log/:roomID", ListLog(service))
}

// ConnectRoom
func ConnectRoom(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := entity.NewUser(c.QueryParam("displayName"))
		roomID, err := entity.StringToID(c.Param("roomID"))
		if err != nil {
			log.Errorf("ListRecentに失敗しました: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		conn, err := generateConnection(c.Response(), c.Request())
		if err != nil {
			log.Errorf("WebSocketの接続に失敗しました: %v", err)
		}
		session := entity.NewSession(user, roomID, conn)
		go writeMessage(session, service)
		go readMessage(session, service)
		return nil
	}
}

// generateConnection
func generateConnection(writer http.ResponseWriter, reader *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return upgrader.Upgrade(writer, reader, nil)
}

// writeMessage
func writeMessage(session *entity.Session, service room.UseCase) {
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		service.DisconnectRoom(session)
		ticker.Stop()
		session.Close()
	}()
	messageChannel, err := service.ConnectRoom(session)
	if err != nil {
		return
	}
	for {
		select {
		case <-session.IsDone:
			return
		case message := <-messageChannel:
			messageResponcePresenter := presenter.MarshalMessage(message)
			if err := session.Conn.WriteJSON(&messageResponcePresenter); err != nil {
				log.Errorf("WebSocketのメッセージの書込に失敗しました: %v", err)
				return
			}
		case <-ticker.C:
			if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("Pingに失敗しました :%v", err)
				return
			}
		}
	}
}

// readMessage
func readMessage(session *entity.Session, service room.UseCase) {
	defer func() {
		service.DisconnectRoom(session)
		session.Close()
	}()
	for {
		if session.IsClosed {
			return
		}
		var messageRequestPresenter presenter.MessageRequestPresenter
		if err := session.Conn.ReadJSON(&messageRequestPresenter); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("WebSocketのメッセージの受信が想定外の原因で失敗しました: %v", err)
			}
			return
		}
		messageRequestPresenter.UserName = session.User.DisplayName.String()
		message := presenter.UnmarshalMessage(&messageRequestPresenter)
		if err := service.Publish(session.RoomID, message); err != nil {
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

// ListLog
func ListLog(service room.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID, err := entity.StringToID(c.Param("roomID"))
		if err != nil {
			log.Errorf("StringToIDに失敗しました: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		messages, err := service.ListRecent(roomID)
		if err != nil {
			log.Errorf("ListRecentに失敗しました: %v", err)
		}
		return c.JSON(http.StatusOK, presenter.MarshalMessages(messages))
	}
}
