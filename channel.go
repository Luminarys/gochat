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
	Ops    []string
}

type Message struct {
	Nick string
	Text string
	isOp bool
}

//Creates and joins a new channel
func NewChannel(channel string, bot *Bot) *Channel {
	Ops := make([]string, 0)
	ready := make(chan bool)

	bot.Conn.AddCallback("353", func(e *irc.Event) {
		for _, nick := range strings.Split(e.Message(), " ") {
			if i := strings.Index(nick, "@"); i == 0 {
				Ops = append(Ops, nick[1:])
			}
		}
		ready <- true
	})

	bot.Conn.Join(channel)
	bot.Conn.SendRaw("NAMES")
	<-ready

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
	for _, mod := range c.Bot.Modules {
		if mod.IsValid(msg) {
			//Handle the action asynchronously
			go func() {
				res := mod.ParseMessage(msg, c)
				if res != "" {
					c.Say(res)
				}
			}()
		}
	}
	c.Buffer = append(c.Buffer, msg)
}
