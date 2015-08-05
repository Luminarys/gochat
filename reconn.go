package gochat

import (
	"errors"
	"net"
)

//Starts Net Listener
func (b *Bot) startNetListener() {
	c := b.Conn

	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		LError.Fatalln("Error, could not listen on the port!")
	}

	conn, err := ln.Accept()
	if err != nil {
		LError.Fatalln("Error, could not accept hijack!")
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

	_, err = conn.Write([]byte(chans))
	if err != nil {
		LWarning.Println("TCP write error: ", err)
		return
	}

	LTrace.Println("Hijacked!")
	c.hijacked = true
	//c.Conn.Close()
	conn.Close()
	ln.Close()
	b.Quit()
}

// Attempt to hijack session previously running bot, returning some basic state information
func hijackSession() (string, error) {
	conn, err := net.Dial("tcp", ":10001")
	if err != nil {
		LWarning.Println("Could not connect to bouncer!")
		return "", err
	}

	var b [4096]byte
	var n int
	n, err = conn.Read(b[:])
	if err != nil {
		LError.Println(err.Error())
		return "", errors.New("Could not receieve channel list!")
	}
	LTrace.Println("Received chan list as: ", string(b[:n]))

	return string(b[:n]), nil
}
