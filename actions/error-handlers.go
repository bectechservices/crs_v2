package actions

import (
	"github.com/gobuffalo/buffalo"
)

func GlobalErrorHandler(status int, err error, c buffalo.Context) error {
	if c.Request().Header.Get("Accept") == "application/json" || c.Request().Header.Get("Content-Type") == "application/json" {
		return c.Error(status, err)
	}
	return c.Render(200, r.HTML("error.html"))
}
