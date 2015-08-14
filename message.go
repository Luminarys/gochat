package gochat

import (
	"errors"
	"strings"
	"time"
)

type Message struct {
	Nick      string
	User      string
	Host      string
	Time      time.Time
	Cmd       string
	Text      string
	Arguments []string
}

//Code shamelessly grabbed from go-ircevent
func ParseMessage(msg string) (m *Message, err error) {
	msg = strings.TrimSuffix(msg, "\n") //Remove \r\n
	msg = strings.TrimSuffix(msg, "\r")
	m = &Message{}
	var source string
	if len(msg) < 5 {
		return nil, errors.New("Malformed msg from server")
	}
	if msg[0] == ':' {
		if i := strings.Index(msg, " "); i > -1 {
			source = msg[1:i]
			msg = msg[i+1 : len(msg)]
		} else {
			return nil, errors.New("Malformed msg from server")
		}

		if i, j := strings.Index(source, "!"), strings.Index(source, "@"); i > -1 && j > -1 && i < j {
			m.Nick = source[0:i]
			m.User = source[i+1 : j]
			m.Host = source[j+1 : len(source)]
		}
	}

	split := strings.SplitN(msg, " :", 2)
	args := strings.Split(split[0], " ")
	m.Cmd = strings.ToUpper(strings.TrimSpace(args[0]))
	m.Arguments = args[1:]
	if len(split) > 1 {
		m.Arguments = append(m.Arguments, split[1])
	}
	if len(m.Arguments) == 0 {
		m.Text = ""
	} else {
		m.Text = m.Arguments[len(m.Arguments)-1]
	}
	m.Time = time.Now()
	return m, nil
}
