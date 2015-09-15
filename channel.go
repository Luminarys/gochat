package gochat

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Representation of a channel or a user query
type Channel struct {
	Name    string
	Buffer  []*Message
	Bot     *Bot
	Ignored map[string]bool
	Ready   bool
	Users   map[string]*User
	Me      *User
	Modules []Module
}

//Creates and joins a new channel
func (bot *Bot) NewChannel(channel string) *Channel {
	users := make(map[string]*User)
	me := &User{}
	ignore := make(map[string]bool)
	ignore[bot.Nick] = true

	c := &Channel{
		Name:    channel,
		Buffer:  make([]*Message, 0),
		Bot:     bot,
		Ignored: ignore,
		Ready:   false,
		Users:   users,
		Me:      me,
	}

	go c.HandleLogs()

	return c
}

func (c *Channel) setUsers(message string) {
	c.Users, c.Me = parseUsers(c.Users, message, c.Name)
	fmt.Println("Updated users list:", c.Users)
	c.Ready = true
}

//Broadcasts a message on a channel.
func (c *Channel) Say(message string) {
	c.Bot.Conn.privmsg(c.Name, message)
	c.Buffer = append(c.Buffer, &Message{Nick: c.Bot.Nick, Text: message})
}

//Leaves a channel and destroys the channel struct
func (c *Channel) Part() {
	c.Bot.Part(c.Name)
}

//Loads a module into the chan
func (c *Channel) AddModule(mod Module) {
	c.Modules = append(c.Modules, mod)
}

//Handles a message in a channel.
func (c *Channel) HandleMessage(msg *Message) {
	fmt.Println(msg.Text)
	c.Buffer = append(c.Buffer, msg)
	//If the nick is not in the ignore list or has their value set to false, then don't process the messages
	if ignored, exists := c.Ignored[msg.Nick]; !ignored || !exists {
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
		//Parse local modules
		for _, mod := range c.Modules {
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
}

//Ignores a nick, preferrably for a bot, but also potentially for spammers
func (c *Channel) IgnoreNick(nick string) {
	c.Ignored[nick] = true
}

//Unignores a nick in the channel
func (c *Channel) UnignoreNick(nick string) {
	c.Ignored[nick] = false
}

//Dumps current logs into a file and wipes the Buffer
func (c *Channel) DumpLogs() {
	if !exists("logs/" + c.Name) {
		os.MkdirAll("logs/"+c.Name, 0755)
	}
	t := time.Now()
	name := strconv.Itoa(t.Year()) + "-" + t.Month().String() + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute())
	f, err := os.Create("logs/" + c.Name + "/" + name)
	if err != nil {
		fmt.Println("Error, could not create log file!")
		return
	}
	defer f.Close()
	for _, msg := range c.Buffer {
		_, err := f.WriteString(msg.Nick + ": " + msg.Text + "\n")
		if err != nil {
			fmt.Println("Error, could not write to log file!")
			return
		}
	}
	c.Buffer = make([]*Message, 0)
}

//Dumps logs every 24 hours to ensure that the buffer doesn't get too big
func (c *Channel) HandleLogs() {
	for {
		time.Sleep(time.Hour * 24)
		c.DumpLogs()
	}
}
