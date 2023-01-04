package entity

import "github.com/google/uuid"

type UID struct {
	value uuid.UUID
}

// NewUID
func NewUID() UID {
	return UID{
		value: uuid.New(),
	}
}

// StringToID
func StringToID(s string) (UID, error) {
	UUID, err := uuid.Parse(s)
	return UID{
		value: UUID,
	}, err
}

// String
func (i UID) String() string {
	return i.value.String()
}

var ErrorUID, _ = StringToID("00000000-0000-0000-0000-000000000000")
