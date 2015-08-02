package gochat

type Channel struct {
	Name   string
	Buffer []*Message
	Bot    *Bot
}

type Message struct {
	Nick string
	Text string
}

//Creates a new channel
func NewChannel(channel string, bot *Bot) *Channel {
	return &Channel{
		Name:   channel,
		Buffer: make([]*Message, 0),
		Bot:    bot,
	}
}

//Says a message on a channel
func (c *Channel) Say(message string) {
	c.Bot.Conn.Privmsg(c.Name, message)
}

//Leaves a channel and destroys the channel struct
func (c *Channel) Part() {
	c.Bot.Conn.Part(c.Name)
	c.Bot.Channels[c.Name] = nil
	c = nil
}

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
}
