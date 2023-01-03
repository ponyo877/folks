package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ponyo877/folks/api/presenter"
)

// MakeRoomHandlers
func MakeRoomHandlers(e *echo.Echo) {
	e.GET("/v1/rooms", ListRoom())
}

// ListRoom
func ListRoom() echo.HandlerFunc {
	return func(c echo.Context) error {
		responce := presenter.Room{
			ID:          "SAMPLE_ID",
			DisplayName: "SAMPLE_DISPLAY_NAME",
			Members: []presenter.User{
				{
					ID:          "SAMPLE_ID",
					DisplayName: "SAMPLE_DISPLAY_NAME",
				},
			},
		}
		return c.JSON(http.StatusOK, responce)
	}
}
