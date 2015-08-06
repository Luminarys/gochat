package gochat

import (
	"testing"
	"time"
)

//Currently broken

func TestHijack(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	time.Sleep(5 * time.Second)
	bot2, err := NewBot("irc.rizon.net:6666", "go-bot", true)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(7 * time.Second)
	bot2.Part(c.Name)
	bot2.Quit()
	time.Sleep(time.Second)
}
