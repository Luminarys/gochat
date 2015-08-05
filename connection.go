package gochat

import (
	"bufio"
	//"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

//Represents a connection to an IRC server
type connection struct {
	Server    string
	ReadChan  chan *Message
	WriteChan chan string
	Conn      net.Conn
	UseTLS    bool
	Nick      string

	//Low level critical modules
	Modules []Module

	// Unix domain socket address for reconnects (linux only)
	unixastr string
	unixlist net.Listener

	// Whether or not this is a reconnect instance
	reconnect bool
	hijacked  bool

	// Duration to wait between sending of messages to avoid being
	// kicked by the server for flooding (default 200ms)
	ThrottleDelay time.Duration
}

//Makes and returns a new connection to a server
func makeConn(server, nick string, UseTLS, recon bool) (*connection, error) {
	conn := new(connection)

	conn.ReadChan = make(chan *Message, 20)
	conn.WriteChan = make(chan string, 20)
	conn.Nick = nick
	conn.unixastr = "/tmp/gochat.sock"
	conn.UseTLS = UseTLS
	conn.ThrottleDelay = time.Millisecond * 200
	conn.hijacked = false
	debug := false
	var err error

	if !debug {
		cmd := exec.Command("go", "run", "bouncer/bouncer.go")
		err = cmd.Start()
		if err != nil {
			LWarning.Println("Could not start bouncer!")
			return nil, err
		}
		time.Sleep(400 * time.Millisecond)
		sConn, err := net.Dial("tcp", "127.0.0.1:10001")
		if err != nil {
			LWarning.Println("Could not connect to bouncer!")
			return nil, err
		}
		_, err = sConn.Write([]byte(server))
		if err != nil {
			LWarning.Println("Could not connect to bouncer!")
			return nil, err
		}
	}
	//Hang until the socket is created
	for !exists("/tmp/gochat.sock") {
	}
	err = conn.connect(server)
	if err != nil {
		LWarning.Println("Could not connect!")
		return nil, err
	}
	LTrace.Println("Connected successfuly!")

	go conn.readMessages()
	go conn.writeMessages()

	return conn, nil
}

func makeReconn(server, nick string, nconn net.Conn) *connection {
	conn := new(connection)

	conn.ReadChan = make(chan *Message, 20)
	conn.WriteChan = make(chan string, 20)
	conn.Nick = nick
	conn.unixastr = fmt.Sprintf("@%s/irc", nick)
	conn.ThrottleDelay = time.Millisecond * 200
	conn.hijacked = true
	conn.Conn = nconn

	go conn.readMessages()
	go conn.writeMessages()

	return conn
}

func (c *connection) user(user string) {
	c.send(fmt.Sprintf("USER %s %d * :%s", user, 8, user))
}

func (c *connection) nick(nick string) {
	c.send("NICK " + nick)
}

func (c *connection) addModule(m Module) {
	c.Modules = append(c.Modules, m)
}

//Establishes connection to a server
func (c *connection) connect(server string) error {
	LTrace.Println("Connecting to server " + server)
	unaddr, err := net.ResolveUnixAddr("unix", c.unixastr)
	if err != nil {
		return errors.New("Could not resolve unix socket")
	}

	c.Conn, err = net.DialUnix("unix", nil, unaddr)
	if err != nil {
		LError.Println(err.Error())
		return errors.New("Could not establish connection, bouncer not running.")
	}
	/*
		if c.UseTLS {
			c.Conn, err = tls.Dial("tcp", server, &tls.Config{})
			if err != nil {
				LWarning.Println("Could not connect to server with TLS")
				return err
			}
		} else {
			c.Conn, err = net.Dial("tcp", server)
			if err != nil {
				LWarning.Println("Could not connect to server")
				return err
			}
		}
	*/
	c.Server = server
	return nil
}

//Quits and destroys connection
func (c *connection) quit() {
	if !c.hijacked {
		c.WriteChan <- "QUIT"
		time.Sleep(time.Millisecond * 50)
		c.Conn.Close()
	}
	c = nil
}

//Sends a private message to a user or channel
func (c *connection) privmsg(who, text string) {
	if !c.hijacked {
		for len(text) > 400 {
			c.send("PRIVMSG " + who + " :" + text[:400])
			text = text[400:]
		}
		c.send("PRIVMSG " + who + " :" + text)
	}
}

func (c *connection) send(msg string) {
	c.WriteChan <- msg
}

//Loop to read messages
func (c *connection) readMessages() {
	LTrace.Println("Started read message loop")
	scan := bufio.NewScanner(c.Conn)
	for scan.Scan() && !c.hijacked {
		msg, err := ParseMessage(scan.Text())
		parsePM(msg)
		if err != nil {
			continue
		}
		used := false
		//Attempt to utilize low level modules, if not then pass it into the chan
		for _, mod := range c.Modules {
			if mod.IsValid(msg, nil) {
				res := mod.ParseMessage(msg, nil)
				if res != "" {
					c.WriteChan <- res
				}
				used = true
				break
			}
		}
		if !used {
			c.ReadChan <- msg
		}
	}
	close(c.ReadChan)
	LTrace.Println("Stopped read message loop")
}

//Loop to write messages
func (c *connection) writeMessages() {
	LTrace.Println("Started write message loop")
	for !c.hijacked {
		select {
		case s := <-c.WriteChan:
			LTrace.Println("Sending Message: " + s)
			_, err := fmt.Fprint(c.Conn, s+"\r\n")
			if err != nil {
				LWarning.Println("Write error, could not send Message("+s+"): ", err.Error())
				return
			}
		}
		//time.Sleep(c.ThrottleDelay)
	}
	close(c.WriteChan)
	LTrace.Println("Stopped write message loop")
}
func parsePM(m *Message) {
	msg := m.Text
	if m.Cmd == "PRIVMSG" && len(msg) > 2 && msg[0] == '\x01' {
		m.Cmd = "CTCP" //Unknown CTCP

		if i := strings.LastIndex(msg, "\x01"); i > 0 {
			msg = msg[1:i]
		} else {
			LWarning.Println("Invalid CTCP Message")
			return
		}

		if msg == "VERSION" {
			m.Cmd = "CTCP_VERSION"

		} else if msg == "TIME" {
			m.Cmd = "CTCP_TIME"

		} else if strings.HasPrefix(msg, "PING") {
			m.Cmd = "CTCP_PING"

		} else if msg == "USERINFO" {
			m.Cmd = "CTCP_USERINFO"

		} else if msg == "CLIENTINFO" {
			m.Cmd = "CTCP_CLIENTINFO"

		} else if strings.HasPrefix(msg, "ACTION") {
			m.Cmd = "CTCP_ACTION"
			if len(msg) > 6 {
				msg = msg[7:]
			} else {
				msg = ""
			}
		}
		m.Arguments[len(m.Arguments)-1] = msg
	}
}
