package http

import (
	"context"
	"fmt"
	"javacode/internal/config"
	"javacode/internal/service"
	"javacode/pkg/logger/sl"
	"log/slog"

	errorslist "javacode/internal/errors"
	middleware "javacode/internal/transport/middleware"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (a *ApiImpl) StartHTTP() error {

	e := echo.New()
	e.POST("/", a.ChangingSumHandler, middleware.CheckHeader)
	err := e.Start(":8180")

	if err != nil {
		return errors.Wrap(err, "ocurred error StartHTTP")
	}
	return nil
}

func (a *ApiImpl) Stop() {

	if err := a.server.Shutdown(context.TODO()); err != nil {
		panic(errors.Wrap(err, "ocurred error Stop"))
	}
}

func (a *ApiImpl) ChangingSumHandler(c echo.Context) error {

	const op = "ChangingSumHandler"
	log := a.log.With(slog.String("op", op))

	trasactionInfo := &service.TrasactionInfo{}

	c.Bind(trasactionInfo)
	if trasactionInfo.Sum == 0 {
		log.Error("Bad Request")
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	trasactionInfo.Role = c.Request().Header.Get(config.HeaderKey)

	log.Info("run ChangingSum")

	sum, err := a.service.ChangingSum(c.Request().Context(), trasactionInfo)
	if err != nil {

		if errors.Is(err, errorslist.ErrInsufficientFunds) {
			log.Error("occurred error Get Sum", sl.Err(err))
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient funds")
		}
		log.Error("occurred error Get Sum", sl.Err(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Info(fmt.Sprintf("changing ok, role - %s , quantity - %d", trasactionInfo.Role, trasactionInfo.Sum))

	resp := &service.RespTrasactionInfo{Sum: sum}
	c.JSON(http.StatusOK, resp)
	return nil
}
