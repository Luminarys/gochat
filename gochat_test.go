package gochat

import (
	"fmt"
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
	time.Sleep(2500 * time.Millisecond)
}

func TestChanJoin(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	time.Sleep(2 * time.Second)
	c.Part()
	bot.Quit()
	time.Sleep(2500 * time.Millisecond)
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("Trying to join chan")
	c := bot.JoinChan("#go-bot-test")
	fmt.Println("Trying to send broadcast")
	bot.Broadcast("broadcast test")
	time.Sleep(2 * time.Second)
	fmt.Println("Trying to send part message")
	c.Part()
	fmt.Println("Trying to send quit message")
	bot.Quit()
	time.Sleep(2500 * time.Millisecond)
}

func TestMessage(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot", false)
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	c.Say("Message test")
	time.Sleep(2 * time.Second)
	c.Part()
	bot.Quit()
	time.Sleep(2500 * time.Millisecond)
}
