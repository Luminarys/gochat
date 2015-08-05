package gochat

import (
	"errors"
	"github.com/mudler/sendfd"
	"net"
)

//All code here is derived from hellabot/whyrusleeping

//Starts Unix Listener
func (b *Bot) startUnixListener() {
	c := b.Conn
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

	fi, err := c.Conn.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}

	err = sendfd.SendFD(con, fi)
	if err != nil {
		panic(err)
	}

	chans := ""
	//Send a comma separated list of channels. Dump logs too
	for chn, _ := range b.Channels {
		chans += chn + ","
	}
	if len(chans) > 2 {
		chans = chans[:len(chans)-1]
	}
	LTrace.Println("Aggregated Chans as string: ", chans)

	_, err = con.Write([]byte(chans))
	if err != nil {
		LWarning.Println("unix listener error: ", err)
		return
	}
	list.Close()

	LTrace.Println("Hijacked!")
	c.hijacked = true
	//c.Conn.Close()
	b.Quit()
}

// Attempt to hijack session previously running bot
func hijackSession(unixAddr string) (net.Conn, string, error) {
	unaddr, err := net.ResolveUnixAddr("unix", unixAddr)
	if err != nil {
		return nil, "", errors.New("Could not resolve unix socket")
	}

	con, err := net.DialUnix("unix", nil, unaddr)
	if err != nil {
		LError.Println(err.Error())
		return nil, "", errors.New("Could not restablish connection, no prior bot.")
	}

	ncon, err := sendfd.RecvFD(con)
	if err != nil {
		LError.Println(err.Error())
		return nil, "", err
	}

	netcon, err := net.FileConn(ncon)
	if err != nil {
		LError.Println(err.Error())
		return nil, "", err
	}

	var b [4096]byte
	var n int
	n, err = con.Read(b[:])
	if err != nil {
		LError.Println(err.Error())
		return nil, "", errors.New("Could not receieve channel list!")
	}
	LTrace.Println("Received chan list as: ", string(b[:n]))

	return netcon, string(b[:n]), nil
}
