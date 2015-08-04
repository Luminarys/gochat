package gochat

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"os"
	"strconv"
	"strings"
	"time"
)

type Channel struct {
	Name    string
	Buffer  []*Message
	Bot     *Bot
	Ops     map[string]bool
	Ignored map[string]bool
}

type Message struct {
	Nick string
	Text string
}

//Creates and joins a new channel
func (bot *Bot) NewChannel(channel string) *Channel {
	Ops := make(map[string]bool, 0)
	ready := make(chan bool)

	i := bot.Conn.AddCallback("353", func(e *irc.Event) {
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
	bot.Conn.RemoveCallback("353", i)

	ignore := make(map[string]bool)
	ignore[bot.Nick] = true

	c := &Channel{
		Name:    channel,
		Buffer:  make([]*Message, 0),
		Bot:     bot,
		Ops:     Ops,
		Ignored: ignore,
	}

	go c.HandleLogs()

	return c
}

//Broadcasts a message on a channel.
func (c *Channel) Say(message string) {
	c.Bot.Conn.Privmsg(c.Name, message)
	c.Buffer = append(c.Buffer, &Message{Nick: c.Bot.Nick, Text: message})
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

//Handles mode changes for users in a chan
func (c *Channel) ModeChange(e *irc.Event) {
	fmt.Println(e.Arguments)
	if e.Arguments[1] == "+o" {
		c.Ops[e.Arguments[2]] = true
	} else if e.Arguments[1] == "-o" {
		c.Ops[e.Arguments[2]] = false
	}
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
