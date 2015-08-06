package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//Bouncer which creates a unix socket to which is used to proxy connections and for session hijacks

//Listens on port 10001 for initial information then creates the unix socket
func main() {
	var err error
	conf_connd := false

	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}

	go func() {
		time.Sleep(1 * time.Second)
		if !conf_connd {
			fmt.Println("No connection within acceptable time limits!")
			os.Exit(0)
		}
	}()

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error, could not accpet connection!")
		os.Exit(1)
	}
	conf_connd = true

	var b [4096]byte
	var n int
	n, err = conn.Read(b[:])
	if err != nil {
		fmt.Println("Could not receieve info!")
	}
	conn.Close()
	ln.Close()
	server := strings.Trim(string(b[:n]), "\r\n")
	fmt.Println("Connecting to server: " + server)

	servConn, err := net.Dial("tcp", "irc.rizon.net:6666")
	if err != nil {
		fmt.Println("Could not connect to server")
		panic(err)
	}
	serverRead := make(chan []byte, 100)
	serverWrite := make(chan []byte, 10)

	go readFromServ(serverRead, servConn)
	go writeToServ(serverWrite, servConn)

	defer servConn.Close()
	//Wait for a client connection, then just copy the two streams into eachother
	l, err := net.Listen("tcp", ":10003")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}
	//Listens for the kill switch
	ld, err := net.Listen("tcp", ":10004")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}
	defer ld.Close()

	var wg sync.WaitGroup
	for {
		cli_connd := false
		go func() {
			time.Sleep(3 * time.Second)
			if !cli_connd {
				servConn.Close()
				l.Close()
				ld.Close()
				os.Exit(0)
			}
		}()
		wg.Add(2)
		clientConn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Receieved client connection!")
		cli_connd = true
		var done = false
		go writeToClient(serverRead, clientConn, &done, &wg)
		go readFromClient(serverWrite, clientConn, &done, &wg)
		tc, err := ld.Accept()
		done = true
		wg.Wait()
		fmt.Println("Receieved client disconn, waiting for reconn")
		clientConn.Close()
		_, err = tc.Write([]byte("ready"))
		if err != nil {
			fmt.Println("TCP write error: ", err)
			return
		}
		tc.Close()
	}
}

func readFromClient(cr chan []byte, sconn net.Conn, done *bool, wg *sync.WaitGroup) {
	to := make(chan bool)
	msg := make(chan []byte, 20)
	for !*done {
		//If msg chan is empty, then execute this, otherwise keep on Timing out
		if len(msg) == 0 {
			//Probably needs to be a better way of handling this
			go func() {
				var b [32 * 1024]byte
				n, err := sconn.Read(b[:])
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
			fmt.Println("Received client message: ", string(m))
			cr <- m
		case <-to:
		}
	}
	wg.Done()
}

func writeToClient(cw chan []byte, cconn net.Conn, done *bool, wg *sync.WaitGroup) {
	to := make(chan bool)
	for !*done {
		go func() {
			time.Sleep(time.Second)
			to <- true
		}()
		select {
		case msg := <-cw:
			fmt.Println("Sending client message: ", string(msg))
			_, err := cconn.Write(msg)
			if err != nil {
				fmt.Println("TCP write error: ", err)
				return
			}
		case <-to:

		}
	}
	wg.Done()
}

func readFromServ(sr chan []byte, sconn net.Conn) {
	var b [32 * 1024]byte
	for {
		n, err := sconn.Read(b[:])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Received server message: ", string(b[:n]))
		sr <- b[:n]
	}
}

func writeToServ(sw chan []byte, sconn net.Conn) {
	for msg := range sw {
		fmt.Println("Sending server message: ", string(msg))
		_, err := sconn.Write(msg)
		if err != nil {
			fmt.Println("TCP write error: ", err)
			return
		}
	}
}
