package gochat

import (
	"github.com/mudler/sendfd"
	"net"
	"time"
)

//All code here is derived from hellabot/whyrusleeping

//Starts Unix Listener
func (c *connection) startUnixListener() {
	unaddr, err := net.ResolveUnixAddr("unix", c.unixastr)
	if err != nil {
		panic(err)
	}
	list, err := net.ListenUnix("unix", unaddr)
	if err != nil {
		panic(err)
	}

	c.unixlist = list
	con, err := list.AcceptUnix()
	if err != nil {
		LWarning.Println("unix listener error: ", err)
		return
	}
	list.Close()

	fi, err := c.Conn.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}

	err = sendfd.SendFD(con, fi)
	if err != nil {
		panic(err)
	}
	LTrace.Println("Hijacked!")
	c.Conn.Close()
	time.Sleep(time.Millisecond * 60)
	select {
	case <-c.ReadChan:
	default:
		close(c.ReadChan)
	}
	close(c.WriteChan)
}

// Attempt to hijack session previously running bot
func (c *connection) hijackSession() bool {
	unaddr, err := net.ResolveUnixAddr("unix", c.unixastr)
	if err != nil {
		LWarning.Println("Could not resolve unix socket")
		return false
	}

	con, err := net.DialUnix("unix", nil, unaddr)
	if err != nil {
		LWarning.Println("Couldnt restablish connection, no prior bot.")
		return false
	}

	ncon, err := sendfd.RecvFD(con)
	if err != nil {
		LError.Println(err.Error())
		return false
	}

	netcon, err := net.FileConn(ncon)
	if err != nil {
		LError.Println(err.Error())
		return false
	}

	c.reconnect = true
	c.Conn = netcon
	return true
}
