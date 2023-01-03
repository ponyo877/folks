package entity

type UID struct {
	Value string
}

func NewUID(value string) UID {
	return UID{
		Value: value,
	}
}

var ErrorUID = NewUID("ERROR")
