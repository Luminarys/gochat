# GoChat
A modular IRC library built in Go

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
The core behind GoChat is the Module interface. Modules must implement a IsValid method which will check whether or not a message should be acted upon, and a ParseMessage function which will parse an input and return an output
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

#TODO
* Move away from the go-ircevent library and become a standalone library
* Database support for logs
* Demo client
