package gochat

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bot, err := NewBot("irc.animebytes.tv:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	bot.Quit()
}

func TestChanJoin(t *testing.T) {
	bot, err := NewBot("irc.animebytes.tv:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Part()
	bot.Quit()
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.animebytes.tv:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	bot.Broadcast("Ayy broadcast")
	c.Part()
	bot.Quit()
}

func TestMessage(t *testing.T) {
	bot, err := NewBot("irc.animebytes.tv:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Say("Ayy message")
	go bot.Conn.Loop()
	time.Sleep(5 * time.Second)
	c.Part()
	bot.Quit()
}
