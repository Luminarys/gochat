package gochat

import (
	"golang.org/x/net/html"
	"net/http"
	"regexp"
)

type Module interface {
	//Takes a message and evaluates whether or not the Module should act upon it
	IsValid(msg *Message, c *Channel) bool
	//Takes a message, and returns the result. If there is no result, "" should be returned
	ParseMessage(msg *Message, c *Channel) string
}

type PingMod struct {
}

func (p *PingMod) IsValid(msg *Message, c *Channel) bool {
	return msg.Text == ".ping"
}

func (p *PingMod) ParseMessage(msg *Message, c *Channel) string {
	return "Pong!"
}

type URLMod struct {
	Re *regexp.Regexp
}

func (u *URLMod) Init() {
	re, err := regexp.Compile("(https?)://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[a-zA-Z0-9+&@#/%=~_|]")
	if err == nil {
		u.Re = re
	}
}

func (u *URLMod) IsValid(msg *Message, c *Channel) bool {
	if u.Re != nil {
		return u.Re.FindString(msg.Text) != ""
	}
	return false
}

func (u *URLMod) ParseMessage(msg *Message, c *Channel) string {
	url := u.Re.FindString(msg.Text)
	response, err := http.Get(url)
	if err != nil {
		return "Error, could not get URL!"
	} else {
		defer response.Body.Close()
		doc, _ := html.Parse(response.Body)
		var f func(*html.Node, bool)
		var re string
		f = func(n *html.Node, pt bool) {
			if pt && n.Type == html.TextNode {
				re = string("Link Title: " + n.Data)
				return
			}
			pt = pt || (n.Type == html.ElementNode && n.Data == "title")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, pt)
			}
		}
		f(doc, false)
		return re
	}
	return ""
}
