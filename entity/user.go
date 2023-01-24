package entity

type User struct {
	ID          UID
	DisplayName DisplayName
}

func NewUser(displayName string) User {
	return User{
		ID:          NewUID(),
		DisplayName: NewDisplayName(displayName),
	}
}
