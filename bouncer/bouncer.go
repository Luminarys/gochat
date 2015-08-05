package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

//Bouncer which creates a unix socket to which is used to proxy connections and for session hijacks

//Listens on port 10001 for initial information then creates the unix socket
func main() {
	var connd = false
	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}

	go func() {
		time.Sleep(1200 * time.Millisecond)
		if !connd {
			ln.Close()
			os.Exit(0)
		}
	}()

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error, could not accpet connection!")
		os.Exit(1)
	}
	connd = true

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

	servConn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Could not connect to server")
		panic(err)
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{"/tmp/gochat.sock", "unix"})
	if err != nil {
		os.Remove("/tmp/gochat.sock")
		panic(err)
	}
	defer os.Remove("/tmp/gochat.sock")
	defer servConn.Close()
	defer l.Close()

	//Wait for a client connection, then just copy the two streams into eachother
	for {
		connd = false
		//Timeout function in case we don't get a connection
		go func() {
			time.Sleep(1200 * time.Millisecond)
			if !connd {
				servConn.Close()
				l.Close()
				os.Remove("/tmp/gochat.sock")
				os.Exit(0)
			}
		}()
		clientConn, err := l.AcceptUnix()
		if err != nil {
			os.Remove("/tmp/gochat.sock")
			panic(err)
		}
		connd = true
		go func() {
			_, err := io.Copy(clientConn, servConn)
			fmt.Println(err)
			if err != nil {
				if err.Error() == "write unix @: broken pipe" {
					servConn.Close()
					l.Close()
					os.Remove("/tmp/gochat.sock")
					os.Exit(0)
				}
			}
			return
		}()
		_, err = io.Copy(servConn, clientConn)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "write unix @: broken pipe" {
				return
			}
		}
	}
}
