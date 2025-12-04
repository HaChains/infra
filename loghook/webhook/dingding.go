package webhook

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ripemd160"

	"github.com/HaChains/infra/loghook"
)

type DingDing struct {
	c   *loghook.Config
	buf *bytes.Buffer
}

// Payload the alarm message.
type Payload struct {
	Content string `json:"content"`
	Group   int    `json:"group"`
	Datakey string `json:"datakey"`
}

// newDingDing new alarm instance.
func newDingDing(c *loghook.Config) *DingDing {
	return &DingDing{
		c:   c,
		buf: new(bytes.Buffer),
	}
}

func (l *DingDing) Send(msg string) error {
	if msg == "" || l.c.Secret == "" {
		return fmt.Errorf("lark alarm send content settings are incorrect")
	}

	body := Payload{
		Content: msg,
		Group:   l.c.Group,
		Datakey: alarmDatakey(msg, l.c.Secret),
	}

	return request(l.c.URL, body, l.buf)
}

func alarmDatakey(str string, version string) string {
	versionNum := []byte(version)
	r := ripemd160.New()
	A := sha256.Sum256([]byte(str))
	V := versionNum
	r.Write(A[:])
	B := append(V, r.Sum(nil)...)
	tempC := sha256.Sum256(B)
	C := sha256.Sum256(tempC[:])
	D := append(B, C[len(C)-8:]...)
	E := hex.EncodeToString(D)
	return E
}
