package repository

import (
	"time"

	"github.com/ponyo877/folks/entity"
)

type RoomMySQLPresenter struct {
	ID          string    `gorm:"column:id;primary_key"`
	DisplayName string    `gorm:"column:display_name"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

type RoomMySQLPresenterList []RoomMySQLPresenter

// TableName
func (p *RoomMySQLPresenter) TableName() string {
	return "rooms"
}

// pickRoom
func (p *RoomMySQLPresenter) pickRoom() (*entity.Room, error) {
	roomID, err := entity.StringToID(p.ID)
	if err != nil {
		return nil, err
	}
	return &entity.Room{
		ID:          roomID,
		DisplayName: entity.NewDisplayName(p.DisplayName),
	}, nil
}

// pickRoomList
func (p *RoomMySQLPresenterList) pickRoomList() ([]*entity.Room, error) {
	var roomList []*entity.Room
	for _, roomMySQLPresenter := range *p {
		room, err := roomMySQLPresenter.pickRoom()
		if err != nil {
			return nil, err
		}
		roomList = append(roomList, room)
	}
	return roomList, nil
}

// roomEntity
func roomEntity(room *entity.Room) *RoomMySQLPresenter {
	return &RoomMySQLPresenter{
		ID:          room.ID.String(),
		DisplayName: room.DisplayName.String(),
	}
}

// CreateRoom
func (r *MessageRepository) CreateRoom(room *entity.Room) error {
	return r.rdb.Create(roomEntity(room)).Error
}

// ListRoom
func (r *MessageRepository) ListRoom() ([]*entity.Room, error) {
	var roomMySQLPresenterList RoomMySQLPresenterList
	if err := r.rdb.Find(&roomMySQLPresenterList).Error; err != nil {
		return nil, err
	}
	return roomMySQLPresenterList.pickRoomList()
}

// GetRoom
func (r *MessageRepository) GetRoom(roomID entity.UID) (*entity.Room, error) {
	var roomMySQLPresenter RoomMySQLPresenter
	if err := r.rdb.First(&roomMySQLPresenter, "id = ?", roomID.String()).Error; err != nil {
		return nil, err
	}
	return roomMySQLPresenter.pickRoom()
}
