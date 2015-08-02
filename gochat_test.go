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

	bot.JoinChan("#go-bot-test")
	bot.Quit()
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.animebytes.tv:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	bot.JoinChan("#go-bot-test")
	bot.Broadcast("Ayy broadcast")
	time.Sleep(2 * time.Second)
	bot.Quit()
}
