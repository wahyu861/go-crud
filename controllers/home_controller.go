package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Welcome to the Home Page!")
}
