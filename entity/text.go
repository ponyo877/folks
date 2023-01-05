package entity

type Text struct {
	value string
}

// StringToText
func StringToText(value string) Text {
	return Text{
		value: value,
	}
}

func (text *Text) String() string {
	return text.value
}
