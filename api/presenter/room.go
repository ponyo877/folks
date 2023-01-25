package presenter

import (
	"github.com/ponyo877/folks/entity"
)

type RoomPresenter struct {
	ID          string            `json:"id"`
	DisplayName string            `json:"name" validate:"required"`
	Users       UserPresenterList `json:"users"`
}

type UserPresenter struct {
	ID          string `json:"id"`
	DisplayName string `json:"name"`
}

type RoomPresenterList []*RoomPresenter

type UserPresenterList []*UserPresenter

// pickUser
func pickUser(user *entity.User) UserPresenter {
	return UserPresenter{
		ID:          user.ID.String(),
		DisplayName: user.DisplayName.String(),
	}
}

// PickUserList
func PickUserList(users []*entity.User) UserPresenterList {
	var userPresenterList UserPresenterList
	for _, user := range users {
		userPresenter := pickUser(user)
		userPresenterList = append(userPresenterList, &userPresenter)
	}
	return userPresenterList
}

// pickRoom
func pickRoom(room *entity.Room) RoomPresenter {
	return RoomPresenter{
		ID:          room.ID.String(),
		DisplayName: room.DisplayName.String(),
		Users:       PickUserList(room.Users),
	}
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
