package presenter

import (
	"github.com/ponyo877/folks/entity"
)

type RoomPresenter struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"name" validate:"required"`
	Members     []UserPresenter `json:"members"`
}

type UserPresenter struct {
	DisplayName string `json:"name"`
}

type RoomPresenterList []*RoomPresenter

// pickRoom
func pickRoom(room *entity.Room) RoomPresenter {
	roomPresenter := RoomPresenter{
		ID:          room.ID.String(),
		DisplayName: room.DisplayName.String(),
	}
	return roomPresenter
}

// PickRoomList
func PickRoomList(roomList []*entity.Room) RoomPresenterList {
	var roomPresenterList RoomPresenterList
	for _, room := range roomList {
		roomPresenter := pickRoom(room)
		roomPresenterList = append(roomPresenterList, &roomPresenter)
	}
	return roomPresenterList
}

// PickUser
func PickUser(user *entity.User) UserPresenter {
	userPresenter := UserPresenter{
		DisplayName: user.DisplayName.String(),
	}
	return userPresenter
}
