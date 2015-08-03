package gcModules

import (
	"github.com/Luminarys/gochat"
)

//Ping module -- Returns Pong! on receiving .ping
type PingMod struct {
}

func (p *PingMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	return msg.Text == ".ping"
}

func (p *PingMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	return "Pong!"
}
