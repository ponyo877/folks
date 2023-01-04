package entity

type Text struct {
	value string
}

// TextToString
func TextToString(value string) Text {
	return Text{
		value: value,
	}
}

func (text *Text) String() string {
	return text.value
}
