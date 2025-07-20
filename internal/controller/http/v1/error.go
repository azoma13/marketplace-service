package v1

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var ErrInvalidGetUserId = errors.New("invalid get userid")

func newErrorResponse(c echo.Context, errStatus int, message string) {
	err := errors.New(message)
	_, ok := err.(*echo.HTTPError)
	if !ok {
		report := echo.NewHTTPError(errStatus, err.Error())
		_ = c.JSON(errStatus, report)
	}
	c.Error(errors.New("internal server error"))
}
