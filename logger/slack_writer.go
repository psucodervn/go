package logger

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type SlackWriter struct {
	webhookURL string
	username   string
	minLevel   zerolog.Level
	client     *resty.Client
	limiter    *rate.Limiter
}

type event struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (s *SlackWriter) Write(p []byte) (n int, err error) {
	var e event
	_ = json.Unmarshal(p, &e)
	_, _ = s.client.R().
		SetBody(map[string]interface{}{
			"text":     fmt.Sprintf("%s: %s\n```%s```", e.Message, e.Error, string(p)),
			"username": s.username,
		}).
		Post(s.webhookURL)
	return len(p), nil
}

func (s *SlackWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < s.minLevel || !s.limiter.Allow() {
		return len(p), nil
	}
	return s.Write(p)
}

func NewSlackWriter(webhookURL, username string, minLevel zerolog.Level) *SlackWriter {
	if _, err := url.Parse(webhookURL); err != nil {
		log.Fatal().Err(err).Str("webhookURL", webhookURL).Msg("invalid webhookURL")
		return nil
	}
	return &SlackWriter{
		username:   username,
		webhookURL: webhookURL,
		minLevel:   minLevel,
		client:     resty.New(),
		limiter: rate.NewLimiter(
			1/10, // 1 event per 10s
			5,
		),
	}
}
