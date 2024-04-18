package middleware

import (
	"javacode/internal/config"
	"net/http"

	"github.com/labstack/echo"
)

func CheckHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch c.Request().Header.Get(config.HeaderKey) {
		case config.HeaderDec:
			return next(c)
		case config.HeaderInc:
			return next(c)
		default:
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}
	}
}
