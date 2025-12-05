package initloghook

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HaChains/infra/loghook"
	"github.com/HaChains/infra/loghook/webhook"
)

type slogHandler struct {
	tag string
	// std log handler
	*slog.TextHandler

	buf     *bytes.Buffer
	bufLock sync.Mutex
	wh      webhook.Webhook
	level   slog.Level
}

func newHookedSlogger(tag string, wh webhook.Webhook) *slogHandler {
	var buf bytes.Buffer

	h := &slogHandler{
		TextHandler: slog.NewTextHandler(os.Stderr, nil),

		buf: &buf,
		wh:  wh,
	}

	if tag != "" {
		h.tag = tag
	}

	return h
}

func (s *slogHandler) setLevel(l string) error {
	switch l {
	case "debug", "DEBUG":
		s.level = slog.LevelDebug
	case "info", "INFO":
		s.level = slog.LevelInfo
	case "warn", "WARN":
		s.level = slog.LevelWarn
	case "error", "ERROR":
		s.level = slog.LevelError
	default:
		return fmt.Errorf("unsupported log level: %s", l)
	}
	return nil
}

func (s *slogHandler) Handle(_ context.Context, r slog.Record) error {
	s.bufLock.Lock()
	s.buf.Reset()
	s.buf.WriteString(r.Level.String())
	s.buf.WriteString(" ")
	s.buf.WriteString(r.Message)
	r.Attrs(func(attr slog.Attr) bool {
		s.buf.WriteString(" ")
		kv := attr.String()
		if strings.Contains(kv, " ") {
			kv = fmt.Sprintf("\"%s\"", kv)
		}
		s.buf.WriteString(kv)
		return true
	})

	msg := s.buf.String()
	s.bufLock.Unlock()

	// send to webhook
	if r.Level >= s.level {
		s.wh.Send(fmt.Sprintf("%s %s", s.tag, msg))
	}

	// write to std
	fmt.Fprintln(os.Stderr, r.Time.Format(time.DateTime), msg)
	return nil
}

func Slog(c *loghook.Config) {
	once.Do(func() {
		wh, err := webhook.NewWebhook(c)
		if err != nil {
			panic(err)
		}

		if wh == nil {
			panic(fmt.Sprintf("webhook init failed, config: %+v", c))
		}
		s := newHookedSlogger(c.Tag, wh)
		if err := s.setLevel(c.Level); err != nil {
			panic(err)
		}
		slog.SetDefault(slog.New(s))
	})
}
