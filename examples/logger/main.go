package main

import (
	"github.com/labstack/echo/v4"
	"github.com/psucodervn/go/logger"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func main() {
	logger.Init(true, false)

	e := echo.New()
	e.HideBanner = true

	e.Use(logger.EchoMiddleware(func(c echo.Context) bool {
		u := strings.ToLower(c.Request().RequestURI)
		return u == "/healthz" || u == "/metrics"
	}))
	e.GET("/", func(c echo.Context) error {
		l := log.Ctx(c.Request().Context())
		l.Info().Msg("index")
		return c.HTML(http.StatusOK, "Hi!")
	})

	log.Err(e.Start(":1234")).Msg("")
}