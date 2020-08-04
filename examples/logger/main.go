package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/psucodervn/go/logger"
)

func init() {
	logger.InitFromEnv()
	log.Error().Msg("init")
}

func main() {
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
