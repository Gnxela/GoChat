package main

import (
	"net"
	"strings"
	"fmt"
)

type User struct {
	connection net.Conn
	queue chan string
}

func (user User) Start() {
	go user.handleConnectionRead();
	go user.handleConnectionWrite();
}

func (user User) Close() {
	user.connection.Close();
}

func (user User) handleConnectionRead() {
	for {
		array := make([]byte, 1024);
		n, err := user.connection.Read(array)
		if(err != nil) {
			if(strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.")) {
				userLeave <- user;
				return
			}else {
				panic(err)
			}
		}
		str := string(array[:n])
		fmt.Println("> " + str)
		
		message := Message{user, str}
		
		userMessage <- message
	}
}

func (user User) handleConnectionWrite() {
	for {
		select {
		case str := <- user.queue:
			array := []byte(str[:len(str) - 1])
			_, err := user.connection.Write(array)
			if(err != nil) {
				panic(err)
			}
		}	
	}
}