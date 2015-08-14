package gcModules

import (
	"github.com/Luminarys/gochat"
	"strings"
)

//Allows for Ignores/Unignores of the specified Nick. Good for dealing with bot feedback loops and spammers. Requires operator status
type IgnoreMod struct {
}

func (m *IgnoreMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	parts := strings.Split(msg.Text, " ")
	if u, exists := c.Users[msg.Nick]; exists && u.CMode >= gochat.Halfop {
		return (parts[0] == ".ignore" || parts[0] == ".unignore") && len(parts) > 1
	}
	return false
}

func (m *IgnoreMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	parts := strings.Split(msg.Text, " ")

	if u, exists := c.Users[parts[1]]; exists && u.CMode >= c.Users[msg.Nick].CMode {
		return msg.Nick + ": Sorry, but I can't ignore or unignore users of mode higher than you"
	}
	if parts[0] == ".ignore" {
		c.IgnoreNick(parts[1])
		return "Ignored " + parts[1]
	} else if parts[0] == ".unignore" {
		c.UnignoreNick(parts[1])
		return "Unignored " + parts[1]
	}
	return ""
}
