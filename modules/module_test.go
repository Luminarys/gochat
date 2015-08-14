package gcModules

import (
	. "github.com/Luminarys/gochat"
	"testing"
	"time"
)

func TestModules(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}

	LoadDefaultModules(bot)

	time.Sleep(1 * time.Second)

	c := bot.JoinChan("#go-bot-test")
	time.Sleep(5 * time.Second)
	c.Part()
	time.Sleep(2 * time.Second)
	bot.Quit()
}
