# GoChat
[![Build Status](https://travis-ci.org/Luminarys/gochat.svg?branch=master)](https://travis-ci.org/Luminarys/gochat) [![GoDoc](https://godoc.org/github.com/Luminarys/gochat?status.png)](https://godoc.org/github.com/Luminarys/gochat) 

A simple modular IRC library written in Go.

# Requirements
Golang, and a working internet connection. An optional dependency for the URL Title module is the golang html package, obtained with `go get golang.org/x/net/html`.

#Setup
A simple bot that joins a test channel then leaves can be written as such:
```
func main() {
    bot, err := gochat.NewBot("irc.rizon.net:6666", "go-bot", "")
    if err != nil {
        //Handle errors
    }
    c := bot.JoinChan("#test-channel")
    c.Part()
    bot.Quit()
}
```
In addition, if you'd like to control verbosity settings, run either `gochat.LogVerbose()`, `gochat.LogWarn()`, or `gochat.LogErr()` prior to creating any bots.
Otherwise, the default logging level will be set to Warning.

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
* Say: The bot will say [msg] on ".say [msg]"

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

#Users
Gochat utilizes a permissions system based on iotas, where users have a mode of:
* Normal
* Voice
* Halfop
* Operator
* Admin
* Owner
These should be used to compare user permissions, e.g. if you want to ensure that a user is operator or higher to execute a command, you would check that `(your channel).Users[(user nick)].CMode >= gochat.Halfop`

#TODO
* Improve connection handling/irc library in general.
* Add in more useful modules
* Use configuration files or flags
* Better persistence and log handling using a proper DB
