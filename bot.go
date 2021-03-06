package gochat

import (
	"errors"
	"os"
	"time"
)

type Bot struct {
	Server   string
	Nick     string
	Channels map[string]*Channel
	Modules  []Module
	Conn     *connection
}

//Creates a new bot for a server, and returns it once it is ready
func NewBot(server, nick, pass string) (*Bot, error) {
	if LTrace == nil || LWarning == nil || LError == nil {
		LogWarn()
	}
	var conn *connection
	var err error

	bot := &Bot{
		Server:   server,
		Nick:     nick,
		Channels: make(map[string]*Channel),
		Modules:  make([]Module, 0),
	}

	conn, err = makeConn(server, nick, false, false) //Create new irc connection
	if err != nil {
		return nil, errors.New("Error! Could not connect")
	}
	bot.Conn = conn
	/*ready once we get the welcome message(001)
	conn.AddCallback("001", func(e *irc.Event) {
		ready <- true
	})
	close(ready)*/

	readyChan := make(chan bool)
	go bot.handleMessages(readyChan)

	bot.Conn.user(nick)
	bot.Conn.nick(nick)
	if pass != "" {
		bot.Auth(pass)
	}

	<-readyChan
	close(readyChan)

	LTrace.Println("Succesfully intialized a bot!")
	return bot, nil
}

func LogVerbose() {
	LogInit(os.Stdout, os.Stderr, os.Stderr)
}

func LogWarn() {
	LogInit(new(NullWriter), os.Stderr, os.Stderr)
}

func LogErr() {
	LogInit(new(NullWriter), new(NullWriter), os.Stderr)
}

func (bot *Bot) handleMessages(ready chan bool) {
	r := false
	LTrace.Println("Starting message handling loop")
	for msg := range bot.Conn.ReadChan {
		//Ignore everything when shutting down
		if bot.Conn.shutdown {
			continue
		}
		// TODO: Split these up into callbacks
		if msg.Cmd == "PRIVMSG" || msg.Cmd == "NOTICE" {
			//Check that channel is valid
			if _, ok := bot.Channels[msg.Arguments[0]]; ok {
				bot.Channels[msg.Arguments[0]].HandleMessage(msg)
			} else {
				//Otherwise it's probably a user, so instantiate a
				//temporary user channel to handle things
				c := bot.NewChannel(msg.Nick)
				c.HandleMessage(msg)
			}
		} else if msg.Cmd == "MODE" || msg.Cmd == "JOIN" || msg.Cmd == "PART" {
			//Requery state information so we don't have to keep track
			//of a bunch of dumb +v +h +o states
			if c, ok := bot.Channels[msg.Arguments[0]]; ok && c != nil {
				c.UpdateUsers()
			}
		} else if msg.Cmd == "353" {
			if c, ok := bot.Channels[msg.Arguments[2]]; ok && c != nil {
				c.setUsers(msg.Text)
			}
		} else if msg.Cmd == "001" {
			if !r {
				ready <- true
				r = true
			}
		}
	}
}

//Loads a module into the bot
func (bot *Bot) AddModule(mod Module) {
	bot.Modules = append(bot.Modules, mod)
}

//Joins a channel, returning the channel once it is ready
func (bot *Bot) JoinChan(chanName string) *Channel {
	LTrace.Println("Creating and registering a channel interally")
	c := bot.NewChannel(chanName)
	bot.Channels[chanName] = c

	readyChan := make(chan bool)
	go func() {
		for !c.Ready {
			//Give other procs a chance to execute
			time.Sleep(100 * time.Millisecond)
		}
		readyChan <- true
	}()

	LTrace.Println("Sending the JOIN message to the server")
	bot.Conn.send("JOIN " + chanName)
	<-readyChan
	LTrace.Println("Channel is ready!")

	return c
}

// Sends a PM to a channel or user
func (bot *Bot) PM(who, text string) {
	bot.Conn.privmsg(who, text)
}

// Sends a NOTICE to a channel or user
func (bot *Bot) Notice(who, text string) {
	bot.Conn.notice(who, text)
}

// Registers the bot
func (bot *Bot) Register(pass string, email string) {
	bot.Conn.register(pass, email)
}

// Sends the password to authenticate
func (bot *Bot) Auth(pass string) {
	bot.Conn.sendPass(pass)
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
