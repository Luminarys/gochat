# GoChat
[![Build Status](https://travis-ci.org/Luminarys/gochat.svg?branch=master)](https://travis-ci.org/Luminarys/gochat) [![GoDoc](https://godoc.org/github.com/Luminarys/gochat?status.png)](https://godoc.org/github.com/Luminarys/gochat) 

A simple modular IRC library written in Go, using go-ircevent.

# Requirements
Gochat currently uses the go-ircevent library, though this is likely to change soon. You can install it by running `go get github.com/thoj/go-ircevent`.

#Setup
A simple bot that joins a test channel then leaves can be written as such:
```
func main() {
    bot, err := gochat.NewBot("irc.rizon.net:6666", "go-bot")
    if err != nil {
        //Handle errors
    }
    c := bot.JoinChan("#test-channel")
    c.Part()
    bot.Quit()
}
```

#Default Modules
These modules are built into the bot and can be easily imported and loaded with the LoadDefaultModules() function. The currently available modules are:
* Ping-Pong: The bot says "Pong!" on a user saying ".ping"
* URL Title Getting: The bot will display the title of a URL whenever a user types on into the chat
* Sedding: The bot will revise your previous message when you say "s/[old phrase]/[new phrase]/"
* Channel Joining: The bot will join a new channel on ".join [channel]". Note that this require operator status to run
* Quoting: The bot will display a somewhat recent quote from a user on ".quote [nick]"
* Cute Pics: The bot will provide a link to a random cute picture obtained from /c/ on ".cute"
* Bot Report: The bot will report in on ".bots"
* Nick Ignore: The bot will ignore or unignore a nick on ".ignore [nick]" and ".unignore [nick]" respectively

To load in the default modules you can write code as such:
```
import (
    //Other imports
    github.com/Luminarys/gochat/modules
)

func main() {
    //Bot initialization above
    gcModules.LoadDefaultModules(bot)
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
* Fix up reconns
* Add in more useful modules
* Allow for global vs. channel modules
* Add in a more comprehensive permissions system
* Use configuration files or flags
* Move away from the go-ircevent library and become a standalone library

