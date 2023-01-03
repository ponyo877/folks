package presenter

type Room struct {
	ID          string `json:"id"`
	DisplayName string `json:"name"`
	Members     []User `json:"members"`
}
