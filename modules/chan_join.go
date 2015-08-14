package gcModules

import (
	"github.com/Luminarys/gochat"
	"strings"
)

//Instructs the bot to join a channel - Requires operator status
type ChanJoinMod struct {
}

func (m *ChanJoinMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	parts := strings.Split(msg.Text, " ")
	cmd := parts[0]
	if cmd != ".join" {
		return false
	} else if u, ok := c.Users[msg.Nick]; ok && u.CMode >= gochat.Operator {
		return true
	}
	return false
}

func (m *ChanJoinMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	parts := strings.Split(msg.Text, " ")
	channel := parts[1]
	c.Bot.JoinChan(channel)
	return "Joined channel!"
}
