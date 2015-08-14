package gcModules

import (
	"github.com/Luminarys/gochat"
    "strings"
)

//Say module. If a user says .say [msg] the bot will repeat into the channel [msg]
type SayMod struct {
}

func (m *SayMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	parts := strings.Split(msg.Text, " ")
	return parts[0] == ".say" && len(parts) == 2
}

func (m *SayMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	parts := strings.Split(msg.Text, " ")
	return parts[1];
}
