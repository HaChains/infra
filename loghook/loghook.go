package loghook

type Config struct {
	Tag         string
	WebHookType string
	URL         string
	Secret      string
	Group       int
	Off         bool
	Level       string
}
