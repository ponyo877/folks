package presenter

type Room struct {
	ID          string `json:"id"`
	DisplayName string `json:"name" validate:"required"`
	Members     []User `json:"members"`
}
