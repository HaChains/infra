package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HaChains/infra/loghook"
	"net/http"
	"time"
)

type Webhook interface {
	Send(message string) error
}

type Type string

const (
	TypeDingDing Type = "Dingding"
	TypeFeishu   Type = "Feishu"
)

func NewWebhook(c *loghook.Config) (Webhook, error) {
	switch c.WebHookType {
	case string(TypeDingDing):
		return newDingDing(c), nil
	case string(TypeFeishu):
		return newFeishu(c), nil
	}
	return nil, fmt.Errorf(`unsupported webhook type: %s`, c.WebHookType)
}

func request(url string, body any, buff *bytes.Buffer) error {
	marshal, err := json.Marshal(body)
	if err != nil {
		return err
	}
	buff.Reset()
	buff.WriteString(string(marshal))

	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark alarm status code %v", resp.StatusCode)
	}

	return nil
}
