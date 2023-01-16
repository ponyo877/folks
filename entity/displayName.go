package entity

type DisplayName struct {
	value string
}

func NewDisplayName(displayNameStr string) DisplayName {
	return DisplayName{
		value: displayNameStr,
	}
}

// String
func (d *DisplayName) String() string {
	return d.value
}
