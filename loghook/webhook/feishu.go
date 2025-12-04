package webhook

import (
	"bytes"

	"github.com/HaChains/infra/loghook"
)

type FeiShu struct {
	c   *loghook.Config
	buf *bytes.Buffer
}

func newFeishu(c *loghook.Config) *FeiShu {
	return &FeiShu{
		c:   c,
		buf: new(bytes.Buffer),
	}
}

func (w *FeiShu) Send(content string) error {
	// construct alarm body
	body := map[string]any{
		"msg_type": "interactive",
		"card": map[string]any{
			"config": map[string]any{
				"wide_screen_mode": true,
			},
			"header": map[string]any{
				"title": map[string]any{
					"tag":     "plain_text",
					"content": "系统报警",
				},
				"template": "red",
			},
			"elements": []map[string]any{
				{
					"tag": "div",
					"text": map[string]any{
						"content": content,
						"tag":     "plain_text",
					},
				},
			},
		},
	}

	return request(w.c.URL, body, w.buf)
}
