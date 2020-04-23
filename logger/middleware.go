package logger

import (
	"context"
	"time"

	"github.com/rs/xid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoLoggerConfig struct {
	Skipper          middleware.Skipper
	HeaderXRequestID string
}

var DefaultEchoLoggerConfig = EchoLoggerConfig{
	Skipper:          middleware.DefaultSkipper,
	HeaderXRequestID: echo.HeaderXRequestID,
}

func EchoMiddleware(skipper middleware.Skipper) echo.MiddlewareFunc {
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if len(rid) == 0 {
				rid = xid.New().String()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, rid)

			l := loggerWithoutCaller.With().Str("request_id", rid).Logger()
			ctx := l.WithContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

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

var ctxIdKey = struct{}{}

func EchoRequestID(skipper middleware.Skipper) echo.MiddlewareFunc {
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			req := c.Request()
			ctx := req.Context()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = xid.New().String()
				ctx = context.WithValue(ctx, ctxIdKey, rid)
				c.SetRequest(req.WithContext(ctx))
			}
			c.Response().Header().Set(echo.HeaderXRequestID, rid)
			return next(c)
		}
	}
}

func EchoRequestLogger(config *EchoLoggerConfig) echo.MiddlewareFunc {
	if config == nil {
		config = new(EchoLoggerConfig)
		*config = DefaultEchoLoggerConfig
	}
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if len(config.HeaderXRequestID) == 0 {
		config.HeaderXRequestID = echo.HeaderXRequestID
	}

	return EchoMiddleware(config.Skipper)
}
