package gochat

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	time.Sleep(time.Second)
	if err != nil {
		t.Error(err.Error())
	}
	bot.Quit()
	time.Sleep(1500 * time.Millisecond)
}

//Currently broken
/*func TestHijack(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	bot.JoinChan("#go-bot-test")
	time.Sleep(5 * time.Second)
	bot2, err := NewBot("irc.rizon.net:6666", "go-bot", true)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(3 * time.Second)
	bot2.Quit()
	time.Sleep(time.Second)
}*/

func TestChanJoin(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	time.Sleep(time.Second)
	c.Part()
	bot.Quit()
	time.Sleep(1500 * time.Millisecond)
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	bot.Broadcast("broadcast test")
	time.Sleep(time.Second)
	c.Part()
	bot.Quit()
	time.Sleep(1500 * time.Millisecond)
}

func TestMessage(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	c.Say("Message test")
	c.Part()
	time.Sleep(time.Second)
	bot.Quit()
	time.Sleep(1500 * time.Millisecond)
}
