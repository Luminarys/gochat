package gcModules

import (
	"github.com/Luminarys/gochat"
	"math/rand"
	"strings"
)

//Returns a random quote from a user
type QuoteMod struct {
}

func (m *QuoteMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	parts := strings.Split(msg.Text, " ")
	return parts[0] == ".quote" && len(parts) > 1
}

func (m *QuoteMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	parts := strings.Split(msg.Text, " ")
	nick := parts[1]
	messages := make([]string, 0)
	for i := len(c.Buffer) - 1; i >= 0; i-- {
		pmsg := c.Buffer[i]
		if pmsg.Nick == nick {
			messages = append(messages, pmsg.Text)
		}
	}
	if len(messages) > 0 {
		return nick + ": " + messages[rand.Intn(len(messages)-1)]
	}
	return "I couldn't find any quotes for that user!"
}
