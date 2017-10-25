package main

import (
	"net"
	"fmt"
	"strings"
	"bytes"
)

var users []User = make([]User, 0)

var userJoin chan User = make(chan User, 0)
var userLeave chan User = make(chan User, 0)
var userMessage chan Message = make(chan Message, 0)

type Message struct {
	sender *User
	message string
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if(err != nil) {
		panic(err)
	}
	go userHandler();
	for {
		connection, err := listener.Accept()
		if(err != nil) {
			panic(err)
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	user := User{connection, make(chan string, 0), "Default"}
	userJoin <- user
}

func userHandler() {
	for {
		select {
		case user := <- userJoin:
			user.Start();
			users = append(users, user)
			fmt.Println("Joined: " + user.connection.RemoteAddr().String())
			sendMessage("Joined: " + user.connection.RemoteAddr().String())
		case user := <- userLeave:
			user.Close();//Maybe not the best idea to do this here, if there is a queue we might try read/write to the connection
			for p, u := range users {
				if(user == u) {
					users = append(users[:p], users[p + 1:]...)
				}
			}
			fmt.Println("Left: " + user.connection.RemoteAddr().String())//Send message after removing user from users
			sendMessage("Left: " + user.connection.RemoteAddr().String())
		case message := <- userMessage:
			if(message.message[0] == '/') {
				command := message.message[1:]
				if(strings.HasPrefix(command, "nick")) {
					message.sender.name = strings.Split(command, " ")[1];
					fmt.Println("Set name to: " + message.sender.name)
				}
			} else {
				var buffer bytes.Buffer
				buffer.WriteString("<")
				buffer.WriteString(message.sender.name)
				buffer.WriteString("> ")
				buffer.WriteString(message.message)
				fmt.Println(message.sender.name)
				sendMessage(buffer.String())//Handle all messages in a single routine so that we ensure that they are ordered correctly for all clients. "correctly" not nessesarily being the right order, but a consistant order
			}
		}
	}
}

func sendMessage(str string) {
	for _, u := range users {
		u.queue <- str;
	}
}