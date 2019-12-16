package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/hlog"
)

func EchoMiddleware(skipper middleware.Skipper) echo.MiddlewareFunc {
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if skipper(c) {
				return next(c)
			}

			l := hlog.FromRequest(c.Request())

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			l.Info().
				Str("ip", c.RealIP()).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Int64("out", res.Size).
				Int64("latency_ms", stop.Sub(start).Milliseconds()).
				Msg("")
			return
		}
	}
}
