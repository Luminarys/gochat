# GoChat
[![Build Status](https://travis-ci.org/Luminarys/GoChat.svg?branch=master)](https://travis-ci.org/Luminarys/GoChat) [![GoDoc](https://godoc.org/github.com/Luminarys/gochat?status.png)](https://godoc.org/github.com/Luminarys/gochat) 

A simple modular IRC library written in Go, using go-ircevent.

# Requirements
Gochat currently uses the go-ircevent library, though this is likely to change soon. You can install it by running `go get github.com/thoj/go-ircevent`.

#Setup
A simple bot that joins a test channel then leaves can be written as such:
```
func main() {
    bot, err := NewBot("irc.rizon.net:6666", "go-bot")
    if err != nil {
        //Handle errors
    }
    c := bot.JoinChan("#test-channel")
    c.Part()
    bot.Quit()
}
```

#Modules
The core behind GoChat is the Module interface. Modules must implement a IsValid function which will check whether or not a message should be acted upon, and a ParseMessage function which will parse an input and return an output.

A simple ping-pong module can be written as such:
```
type PingMod struct {
}

func (p *PingMod) IsValid(msg *Message, c *Channel) bool {
    //Returns true if the message is equal to ".ping" and false otherwise
    return msg.Text == ".ping"
}

func (p *PingMod) ParseMessage(msg *Message, c *Channel) string {
    //Return "Pong!" to be printed
    return "Pong!"
}

func main() {
    //Declare bot stuff
    bot.AddModule(&{PingMod}) 
}

```
Now, whenever a user types ".ping" into a channel, the bot will respond with "Pong!".

#TODO
* Add in more useful modules
* Allow for global vs. channel modules
* Add in a more comprehensive permissions system
* Use configuration files or flags
* Move away from the go-ircevent library and become a standalone library

