package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/psucodervn/go/logger"
	"github.com/psucodervn/go/validator"
)

func NewDefaultEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = validator.NewStructValidator()

	e.Use(middleware.Recover())
	e.Use(logger.EchoMiddleware(func(c echo.Context) bool {
		uri := c.Request().RequestURI
		return uri == "/healthz" || uri == "/metrics"
	}))
	e.HTTPErrorHandler = ErrorHandler(e)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return e
}

func ErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		if errs, ok := err.(validator.Errors); ok {
			code = http.StatusUnprocessableEntity
			msg = echo.Map{"message": "validation failed: " + errs.Error(), "errors": errs}
		} else if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message
			if he.Internal != nil {
				err = fmt.Errorf("%v, %v", err, he.Internal)
			}
		} else if e.Debug {
			msg = err.Error()
		} else {
			msg = http.StatusText(code)
		}
		if _, ok := msg.(string); ok {
			msg = echo.Map{"success": false, "message": msg}
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}
			if err != nil {
				log.Ctx(c.Request().Context()).Err(err).Msg("")
			}
		}
	}
}
