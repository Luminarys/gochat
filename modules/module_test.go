package gcModules

import (
	. "github.com/Luminarys/gochat"
	"testing"
	"time"
)

func TestModules(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	LoadDefaultModules(bot)

	c := bot.JoinChan("#go-bot-test")
	time.Sleep(18 * time.Second)
	c.Part()
	bot.Quit()
}
