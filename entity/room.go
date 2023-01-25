package entity

type Room struct {
	ID          UID
	DisplayName DisplayName
	Users       []*User
}

func NewRoom(displayName DisplayName) *Room {
	return &Room{
		ID:          NewUID(),
		DisplayName: displayName,
	}
}

// SetUsers
func (r *Room) AddUsers(users []*User) *Room {
	r.Users = users
	return r
}
