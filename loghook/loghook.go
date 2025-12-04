package loghook

type Config struct {
	WebHookType string
	URL         string
	Secret      string
	Group       int
	Off         bool
	Level       string
}
