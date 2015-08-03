package gochat

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	bot.Quit()
}

func TestChanJoin(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Part()
	bot.Quit()
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	bot.Broadcast("broadcast test")
	c.Part()
	bot.Quit()
}

func TestMessage(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Say("Message test")
	c.Part()
	bot.Quit()
}

func TestModules(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	LoadDefaultModules(bot)

	c := bot.JoinChan("#go-bot-test")
	time.Sleep(10 * time.Second)
	c.DumpLogs()
	c.Part()
	bot.Quit()
}

//Loads the default provided Modules
func LoadDefaultModules(bot.Bot) {
	bot.AddModule(&PingMod{})
	um := &URLMod{}
	um.Init()
	bot.AddModule(um)
	sm := &SedMod{}
	sm.Init()
	bot.AddModule(sm)
	bot.AddModule(&QuoteMod{})
	qm := &CuteMod{}
	qm.Init()
	bot.AddModule(qm)
	bot.AddModule(&ChanJoinMod{})
}
