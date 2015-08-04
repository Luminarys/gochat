package gochat

import (
	"errors"
	"os"
)

type Bot struct {
	Server   string
	Nick     string
	Channels map[string]*Channel
	Modules  []Module
	Conn     *connection
}

//Creates a new bot for a server, and returns it once it is ready
func NewBot(server, nick string, hijack bool) (*Bot, error) {
	logInit(os.Stdout, os.Stderr, os.Stderr)
	conn, err := makeConn(server, nick, false, hijack) //Create new irc connection
	if err != nil {
		return nil, errors.New("Error! Could not connect")
	}

	/*ready once we get the welcome message(001)
	conn.AddCallback("001", func(e *irc.Event) {
		ready <- true
	})
	close(ready)*/

	bot := &Bot{
		Server:   server,
		Nick:     nick,
		Channels: make(map[string]*Channel),
		Modules:  make([]Module, 0),
		Conn:     conn,
	}

	ready := make(chan bool)
	go bot.handleMessages(ready)

	bot.Conn.user(nick)
	bot.Conn.nick(nick)

	<-ready
	close(ready)

	return bot, nil
}

func (bot *Bot) handleMessages(ready chan bool) {
	r := false
	for msg := range bot.Conn.ReadChan {
		if !r {
			ready <- true
			r = true
		}
		LTrace.Println(msg.Cmd+": ", msg.Text)
		if msg.Cmd == "PRIVMSG" {
			bot.Channels[msg.Arguments[0]].HandleMessage(msg)
		} else if msg.Cmd == "MODE" {
			bot.Channels[msg.Arguments[0]].ModeChange(msg)
		} else if msg.Cmd == "353" {
			bot.Channels[msg.Arguments[2]].SetOps(msg.Text)
		}
	}
}

//Loads a module into the bot
func (bot *Bot) AddModule(mod Module) {
	bot.Modules = append(bot.Modules, mod)
}

//Joins a channel
func (bot *Bot) JoinChan(chanName string) *Channel {
	c := bot.NewChannel(chanName)
	bot.Channels[chanName] = c
	bot.Conn.send("JOIN " + chanName)
	bot.Conn.send("NAMES " + chanName)
	return c
}

//Disconnects and destroys the bot
func (bot *Bot) Quit() {
	bot.Conn.quit()
}

//Broadcasts a message to all chans
func (bot *Bot) Broadcast(message string) {
	for _, c := range bot.Channels {
		bot.Conn.privmsg(c.Name, message)
	}
}

//Leaves a channel and destroys it
func (bot *Bot) Part(channel string) {
	bot.Conn.send("PART " + channel)
	bot.Channels[channel] = nil
}
