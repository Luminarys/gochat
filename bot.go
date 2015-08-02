package gochat

import (
	"errors"
	"github.com/thoj/go-ircevent"
)

type Bot struct {
	Server   string
	Nick     string
	Channels []*Channel
	Conn     *irc.Connection
}

//Creates a new bot for a server, and returns it once it is ready
func NewBot(server string, nick string) (*Bot, error) {
	bot := irc.IRC(nick, nick) //Create new ircobj
	if err := bot.Connect(server); err != nil {
		return nil, errors.New("Error! Could not connect")
	}
	ready := make(chan bool)

	//ready once we get the welcome message(001)
	bot.AddCallback("001", func(e *irc.Event) {
		ready <- true
	})
	<-ready

	return &Bot{
		Server:   server,
		Nick:     nick,
		Channels: make([]*Channel, 0),
		Conn:     bot,
	}, nil
}

//Joins a channel
func (bot *Bot) JoinChan(chanName string) error {
	c := newChannel(chanName)
	bot.Conn.Join(chanName)
	bot.Channels = append(bot.Channels, c)
	return nil
}

//Disconnects and destroys the bot
func (bot *Bot) Quit() {
	bot.Conn.Quit()
	bot = nil
}

//Broadcasts a message to all chans
func (bot *Bot) Broadcast(message string) {
	for _, c := range bot.Channels {
		bot.Conn.Privmsg(c.Name, message)
	}
}
