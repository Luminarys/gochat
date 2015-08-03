package gcModules

import (
	"github.com/Luminarys/gochat"
	"regexp"
	"strings"
)

//Allows user to perform a s/[original]/[new]/ on previous text
type SedMod struct {
	Re *regexp.Regexp
}

func (s *SedMod) Init() {
	re, err := regexp.Compile("s/[-a-zA-Z0-9+&@#/%?=~_|!:,.;]+/[-a-zA-Z0-9+&@#/%?=~_|!:,.;]+/")
	if err == nil {
		s.Re = re
	}
}

func (s *SedMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	if s.Re != nil {
		return s.Re.MatchString(msg.Text)
	}
	return false
}

func (s *SedMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	sub := s.Re.FindString(msg.Text)
	subArr := strings.Split(sub, "/")
	prev := false
	for i := len(c.Buffer) - 1; i >= 0 && i >= len(c.Buffer)-100; i-- {
		pmsg := c.Buffer[i]
		if pmsg.Nick == msg.Nick {
			if !prev {
				prev = true
				continue
			}
			return msg.Nick + ": " + strings.Replace(pmsg.Text, subArr[1], subArr[2], -1)
		}
	}
	return "Luminarys: I couldn't find any recent messages to perform a subsitution for!"
}
