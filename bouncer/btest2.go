package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

//Bouncer which creates a unix socket to which is used to proxy connections and for session hijacks

//Listens on port 10001 for initial information then creates the unix socket
func main() {
	ln, err := net.Listen("tcp", ":10001")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}
	fmt.Println("Waiting for init client conn")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error, could not accpet connection!")
		os.Exit(1)
	}

	fmt.Println("Got init conn")
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
	defer servConn.Close()
	//servWriter := io.MultiWriter(servConn, os.Stdout)

	l, err := net.Listen("tcp", ":10003")
	if err != nil {
		fmt.Println("Error, could not listen on the port!")
		os.Exit(1)
	}
	defer l.Close()

	//Wait for a client connection, then just copy the two streams into eachother
	for {
		fmt.Println("Waiting for client connection")
		done := false
		d1 := false
		//Timeout function in case we don't get a connection
		clientConn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Receieved a client connection")
		//clientWriter := io.MultiWriter(clientConn, os.Stdout)
		go func() {
			_, err := io.Copy(clientConn, servConn)
			fmt.Println(err)
			if err != nil {
				if !done {
					servConn.Close()
					l.Close()
					os.Exit(0)
				}
			} else {
				done = true
			}
			d1 = true
			return
		}()
		_, err = io.Copy(servConn, clientConn)
		if err != nil {
			fmt.Println(err)
			if !done {
				return
			}
		} else {
			done = true
		}
		fmt.Println("Waiting for both goroutines to quit")
		for !d1 {
		}
		fmt.Println("Client DCd succesfully")
	}
}
