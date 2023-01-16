package repository

import (
	"github.com/ponyo877/folks/entity"
)

// CreateRoom
func (r *MessageRepository) CreateRoom(displayName entity.DisplayName) error {
	// newRoomID := entity.NewUID().String()
	// return r.kvs.SetNX(context.Background(), "room:"+newRoomID, displayName.String(), 10*time.Second).Err()
	return nil
}

// ListRoom
func (r *MessageRepository) ListRoom() ([]*entity.Room, error) {
	// roomStrList, _, err := r.kvs.Scan(context.Background(), 0, "room:*", 0).Result()
	// if err != nil {
	// 	return []*entity.Room{}, err
	// }
	// var rooms []*entity.Room
	// for _, roomStr := range roomStrList {
	// 	rooms = append(rooms, entity.NewRoom(en))
	// }
	return []*entity.Room{}, nil
}
