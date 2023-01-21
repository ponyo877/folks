package entity

type User struct {
	DisplayName DisplayName
}

func NewUser(displayName string) User {
	return User{
		DisplayName: NewDisplayName(displayName),
	}
}
