package gochat

import (
	"errors"
	//"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
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
	// Whther or not the connection is being hijacked
	hijacked bool
	shutdown bool

	wg *sync.WaitGroup

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
	conn.UseTLS = UseTLS
	conn.ThrottleDelay = time.Millisecond * 200
	conn.hijacked = false
	conn.shutdown = false
	var wg sync.WaitGroup
	//Require that the read loop, write loop, and hijack listener be shutdown
	conn.wg = &wg
	var err error
	debug := false

	if !recon && !debug {
		cmd := exec.Command("go", "run", "bouncer/bouncer.go")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			LWarning.Println("Could not get stdout bouncer!")
			return nil, err
		}
		go io.Copy(os.Stdout, stdout)
		err = cmd.Start()
		if err != nil {
			LWarning.Println("Could not start bouncer!")
			return nil, err
		}
		time.Sleep(400 * time.Millisecond)
		sConn, err := net.Dial("tcp", ":10001")
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
	time.Sleep(500 * time.Millisecond)
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
	var err error

	conn, err := net.Dial("tcp", ":10003")
	if err != nil {
		LError.Println(err.Error())
		return errors.New("Could not establish connection, bouncer not running.")
	}
	c.Conn = conn

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
		c.wg.Add(3)
		c.WriteChan <- "QUIT"
		time.Sleep(50 * time.Millisecond)
	} else {
		c.wg.Add(2)
	}
	c.shutdown = true
	c.wg.Wait()
	c.Conn.Close()
}

//Sends a private message to a user or channel
func (c *connection) privmsg(who, text string) {
	for len(text) > 400 {
		c.send("PRIVMSG " + who + " :" + text[:400])
		text = text[400:]
	}
	c.send("PRIVMSG " + who + " :" + text)
}

func (c *connection) send(msg string) {
	c.WriteChan <- msg
}

//Loop to read messages
func (c *connection) readMessages() {
	LTrace.Println("Started read message loop")
	to := make(chan bool)
	msg := make(chan []byte, 20)
	for !c.shutdown {
		//If msg chan is empty, then execute this, otherwise keep on Timing out
		if len(msg) == 0 {
			//Probably needs to be a better way of handling this
			go func() {
				var b [32 * 1024]byte
				n, err := c.Conn.Read(b[:])
				if err != nil {
					fmt.Println("Error reading from client: ", err.Error())
					return
				}
				msg <- b[:n]
			}()
		}
		go func() {
			time.Sleep(time.Second)
			to <- true
		}()
		select {
		case m := <-msg:
			fmt.Println("Received server message: ", string(m))
			msg, err := ParseMessage(string(m))
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
		case <-to:
			LTrace.Println("Read loop timed out!")
		}
	}
	close(c.ReadChan)
	LTrace.Println("Stopped read message loop")
	c.wg.Done()
}

//Loop to write messages
func (c *connection) writeMessages() {
	LTrace.Println("Started write message loop")
	to := make(chan bool)
	for !c.shutdown {
		go func() {
			time.Sleep(time.Second)
			to <- true
		}()
		select {
		case msg := <-c.WriteChan:
			fmt.Println("Sending client message: ", string(msg))
			_, err := c.Conn.Write([]byte(msg + "\r\n"))
			if err != nil {
				LWarning.Println("Write error, could not send Message("+msg+"): ", err.Error())
				return
			}
		case <-to:
			LTrace.Println("Write timeout occured")
		}
	}
	close(c.WriteChan)
	LTrace.Println("Stopped write message loop")
	c.wg.Done()
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
