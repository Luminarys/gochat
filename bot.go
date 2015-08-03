package gochat

import (
	"errors"
	"github.com/thoj/go-ircevent"
)

type Bot struct {
	Server   string
	Nick     string
	Channels map[string]*Channel
	Modules  []Module
	Conn     *irc.Connection
}

//Creates a new bot for a server, and returns it once it is ready
func NewBot(server string, nick string) (*Bot, error) {
	conn := irc.IRC(nick, nick) //Create new irc connection
	if err := conn.Connect(server); err != nil {
		return nil, errors.New("Error! Could not connect")
	}

	ready := make(chan bool)
	//ready once we get the welcome message(001)
	conn.AddCallback("001", func(e *irc.Event) {
		ready <- true
	})
	<-ready
	close(ready)

	bot := &Bot{
		Server:   server,
		Nick:     nick,
		Channels: make(map[string]*Channel),
		Modules:  make([]Module, 0),
		Conn:     conn,
	}

	//Whenever a message is detected, send it to the respective channel for handling
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		bot.Channels[e.Arguments[0]].HandleMessage(&Message{Nick: e.Nick, Text: e.Message()})
	})

	//Whenever a message is detected, send it to the respective channel for handling
	conn.AddCallback("MODE", func(e *irc.Event) {
		if e.Nick != nick {
			bot.Channels[e.Arguments[0]].ModeChange(e)
		}
	})

	return bot, nil
}

//Loads the default provided Modules
func (bot *Bot) LoadDefaultModules() {
	bot.AddModule(&PingMod{})
	um := &URLMod{}
	um.Init()
	bot.AddModule(um)
	sm := &SedMod{}
	sm.Init()
	bot.AddModule(sm)
}

//Loads a module into the bot
func (bot *Bot) AddModule(mod Module) {
	bot.Modules = append(bot.Modules, mod)
}

//Joins a channel
func (bot *Bot) JoinChan(chanName string) *Channel {
	c := NewChannel(chanName, bot)
	bot.Channels[chanName] = c
	return c
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
