package gochat

import (
	"encoding/json"
	"golang.org/x/net/html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Module interface {
	//Takes a message and evaluates whether or not the Module should act upon it
	IsValid(msg *Message, c *Channel) bool
	//Takes a message, and returns the result. If there is no result, "" should be returned
	ParseMessage(msg *Message, c *Channel) string
}

//Ping module -- Returns Pong! on receiving .ping
type PingMod struct {
}

func (p *PingMod) IsValid(msg *Message, c *Channel) bool {
	return msg.Text == ".ping"
}

func (p *PingMod) ParseMessage(msg *Message, c *Channel) string {
	return "Pong!"
}

//URL title detection module -- Returns the URL title when it detects a URL
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

func (s *SedMod) IsValid(msg *Message, c *Channel) bool {
	if s.Re != nil {
		return s.Re.MatchString(msg.Text)
	}
	return false
}

func (s *SedMod) ParseMessage(msg *Message, c *Channel) string {
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

//Returns a random quote from a user
type QuoteMod struct {
}

func (m *QuoteMod) IsValid(msg *Message, c *Channel) bool {
	parts := strings.Split(msg.Text, " ")
	return parts[0] == ".quote" && len(parts) > 1
}

func (m *QuoteMod) ParseMessage(msg *Message, c *Channel) string {
	parts := strings.Split(msg.Text, " ")
	nick := parts[1]
	messages := make([]string, 0)
	for i := len(c.Buffer) - 1; i >= 0; i-- {
		pmsg := c.Buffer[i]
		if pmsg.Nick == nick {
			messages = append(messages, pmsg.Text)
		}
	}
	if len(messages) > 0 {
		return nick + ": " + messages[rand.Intn(len(messages)-1)]
	}
	return "I couldn't find any quotes for that user!"
}

//Returns a link to a random cute pic, courtesy of /c/
type CuteMod struct {
}

type Catalog []Page

type Page struct {
	Number  int             `json:"page"`
	Threads []*FourchanPost `json:"threads"`
}

type Thread struct {
	Posts []*FourchanPost `json:"posts"`
}

type FourchanPost struct {
	Number        int    `json:"no"`
	ReplyTo       int    `json:"resto"`
	Sticky        int    `json:"sticky,omitempty"`
	Closed        int    `json:"closed,omitempty"`
	Date          string `json:"now"`
	Timestamp     int    `json:"time"`
	Name          string `json:"name,omitempty"` //username
	Tripcode      string `json:"trip,omitempty"`
	ID            string `json:"id,omitempty"` //user ID
	Capcode       string `json:"capcode,omitempty"`
	CountryName   string `json:"country_name,omitempty"`
	Email         string `json:"email,omitempty"`
	Subject       string `json:"sub,omitempty"`
	Text          string `json:"com,omitempty"` //HTML
	FileTime      uint64 `json:"tim,omitempty"`
	FileExt       string `json:"ext,omitempty"`
	FileDeleted   int    `json:"filedeleted,omitempty"`
	Spoiler       int    `json:"spoiler,omitempty"`
	OmitedPosts   int    `json:"omitted_posts,omitempty"`
	OmittedImages int    `json:"omitted_images,omitempty"`
	Replies       int    `json:"replies,omitempty"`
	Images        int    `json:"images,omitempty"`
}

func (m *CuteMod) IsValid(msg *Message, c *Channel) bool {
	return msg.Text == ".cute"
}

func (m *CuteMod) ParseMessage(msg *Message, c *Channel) string {
	response, err := http.Get("https://a.4cdn.org/c/catalog.json")
	urls := make([]string, 0)
	if err != nil {
		return "Error, could not get URL!"
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		var data Catalog
		json.Unmarshal(body, &data)
		for _, page := range data {
			for _, thread := range page.Threads {
				r, err := http.Get("https://a.4cdn.org/c/thread/" + strconv.Itoa(thread.Number) + ".json")
				if err == nil {
					b, err := ioutil.ReadAll(r.Body)
					if err != nil {
						continue
					}
					var t Thread
					err = json.Unmarshal(b, &t)
					if err != nil {
						continue
					}
					for _, post := range t.Posts {
						if post.FileTime != 0 {
							urls = append(urls, "https://i.4cdn.org/c/"+strconv.FormatUint(post.FileTime, 10)+post.FileExt)
						}
					}
					r.Body.Close()
				} else {
					continue
				}
			}
		}
	}
	if len(urls) > 0 {
		return "Here's a random cute pic: " + urls[rand.Intn(len(urls)-1)]
	}
	return "I couldn't find anything cute, you may want to try again later!"
}
