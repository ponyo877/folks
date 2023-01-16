package entity

type Room struct {
	ID          UID
	DisplayName DisplayName
}

func NewRoom(displayName DisplayName) Room {
	return Room{
		ID:          NewUID(),
		DisplayName: displayName,
	}
}
