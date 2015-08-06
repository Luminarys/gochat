package gochat

import (
	"errors"
	"net"
	"time"
)

//Starts Net Listener
func (b *Bot) startNetListener() {
	c := b.Conn

	ln, err := net.Listen("tcp", ":10002")
	if err != nil {
		LError.Println("Error, could not listen on the port: " + err.Error())
		return
	}
	//Shutdown the net socket if shutdown is receieved and we're not being hijacked
	go func() {
		for !b.Conn.shutdown {

		}
		ln.Close()
	}()

	conn, err := ln.Accept()
	if err != nil && !b.Conn.shutdown {
		LError.Fatalln("Error, could not accept hijack: " + err.Error())
	} else if b.Conn.shutdown {
		c.wg.Done()
		return
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

	//Shut everything down
	c.wg.Add(2)
	c.shutdown = true
	c.wg.Wait()
	//c.Conn.Close()

	_, err = conn.Write([]byte("Done!"))
	if err != nil {
		LWarning.Println("TCP write error: ", err)
		return
	}
	LTrace.Println("Sent done message, closing up")
	c.Conn.Close()
	conn.Close()
	ln.Close()
}

// Attempt to hijack session previously running bot, returning some basic state information
func hijackSession() (string, error) {
	conn, err := net.Dial("tcp", ":10002")
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

	n, err = conn.Read(b[:])
	if err != nil {
		LError.Println(err.Error())
		return "", errors.New("Could not receieve done signal!")
	}
	LTrace.Println("Previous bot is done!")

	//Sends kill message to the server
	tc, err := net.Dial("tcp", ":10004")
	if err != nil {
		LWarning.Println("Could not send kill message to server!")
		time.Sleep(5 * time.Second)
		return "", err
	}
	n, err = tc.Read(b[:])
	if err != nil {
		LError.Println(err.Error())
		return "", errors.New("Could not receieve bouncer ready signal!")
	}
	LTrace.Println("Bouncer is ready!")
	tc.Close()
	time.Sleep(2 * time.Second)
	conn.Close()
	return string(b[:n]), nil
}
