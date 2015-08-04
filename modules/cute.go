package gcModules

import (
	"encoding/json"
	"github.com/Luminarys/gochat"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//Returns a link to a random cute pic, courtesy of /c/
type CuteMod struct {
	Updating bool
	Pics     []string
}

type catalog []page

type page struct {
	Number  int             `json:"page"`
	Threads []*fourchanPost `json:"threads"`
}

type bread struct {
	Posts []*fourchanPost `json:"posts"`
}

type fourchanPost struct {
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

func (m *CuteMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	return msg.Text == ".cute"
}

func (m *CuteMod) Update() {
	m.Updating = true
	response, err := http.Get("https://a.4cdn.org/c/catalog.json")
	urls := make([]string, 0)
	if err != nil {
		return
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		var data catalog
		json.Unmarshal(body, &data)
		for _, page := range data {
			for _, thread := range page.Threads {
				r, err := http.Get("https://a.4cdn.org/c/thread/" + strconv.Itoa(thread.Number) + ".json")
				if err == nil {
					b, err := ioutil.ReadAll(r.Body)
					if err != nil {
						continue
					}
					var t bread
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
	m.Pics = urls
	m.Updating = false
}

func (m *CuteMod) Init() {
	go m.Update()
	go func() {
		time.Sleep(30 * time.Minute)
		m.Update()
	}()
}

func (m *CuteMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	if !m.Updating {
		if len(m.Pics) > 0 {
			return "Here's a random cute pic: " + m.Pics[rand.Intn(len(m.Pics)-1)]
		}
		return "I couldn't find anything cute, you may want to try again later!"
	} else {
		return "I'm currently busy refilling with moe, try again later!!"
	}
}
