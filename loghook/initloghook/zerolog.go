package initloghook

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/HaChains/infra/loghook"
	"github.com/HaChains/infra/loghook/webhook"
)

func ZeroLog(c *loghook.Config) {
	once.Do(func() {
		wh, err := webhook.NewWebhook(c)
		if err != nil {
			panic(err)
		}
		hook := &zl{
			c:  c,
			wh: wh,
		}
		err = hook.level.UnmarshalText([]byte(c.Level))
		if err != nil {
			panic(fmt.Errorf("unmarshal log level error: %v", err))
		}
		log.Hook(hook)
	})
}

type zl struct {
	c     *loghook.Config
	level zerolog.Level

	wh webhook.Webhook
}

func (z *zl) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if z.c.Off == true || level < z.level {
		return
	}

	err := z.wh.Send(message)
	if err != nil {
		log.Error().Msgf("zerolog hook: webhook send error: %v", err)
	}
}
