package gochat

type Channel struct {
	Name   string
	Buffer []string
}

func newChannel(channel string) *Channel {
	return &Channel{
		Name:   channel,
		Buffer: make([]string, 0),
	}
}
