package gochat

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

type Channel struct {
	Name   string
	Buffer []*Message
	Bot    *Bot
	Ops    map[string]bool
}

type Message struct {
	Nick string
	Text string
}

//Creates and joins a new channel
func NewChannel(channel string, bot *Bot) *Channel {
	Ops := make(map[string]bool, 0)
	ready := make(chan bool)

	bot.Conn.AddCallback("353", func(e *irc.Event) {
		for _, nick := range strings.Split(e.Message(), " ") {
			if i := strings.Index(nick, "@"); i == 0 {
				Ops[nick[1:]] = true
			} else {
				//Have to deal with other signs like %
				Ops[nick] = false
			}
		}
		ready <- true
	})

	bot.Conn.Join(channel)
	bot.Conn.SendRaw("NAMES")
	<-ready
	close(ready)

	return &Channel{
		Name:   channel,
		Buffer: make([]*Message, 0),
		Bot:    bot,
		Ops:    Ops,
	}
}

//Broadcasts a message on a channel.
func (c *Channel) Say(message string) {
	c.Bot.Conn.Privmsg(c.Name, message)
}

//Leaves a channel and destroys the channel struct
func (c *Channel) Part() {
	c.Bot.Conn.Part(c.Name)
	c.Bot.Channels[c.Name] = nil
	c = nil
}

//Handles a message in a channel.
func (c *Channel) HandleMessage(msg *Message) {
	fmt.Println(msg.Text)
	c.Buffer = append(c.Buffer, msg)
	for _, mod := range c.Bot.Modules {
		if mod.IsValid(msg, c) {
			//Handle the action asynchronously
			go func(mod Module) {
				res := mod.ParseMessage(msg, c)
				if res != "" {
					c.Say(res)
				}
			}(mod)
		}
	}
}

func (c *Channel) ModeChange(e *irc.Event) {
	fmt.Println(e.Arguments)
	if e.Arguments[1] == "+o" {
		c.Ops[e.Arguments[2]] = true
	} else if e.Arguments[1] == "-o" {
		c.Ops[e.Arguments[2]] = false
	}
}
