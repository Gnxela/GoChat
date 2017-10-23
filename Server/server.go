package main

import (
	"net"
	"fmt"
)

var users []User = make([]User, 0)

var userJoin chan User = make(chan User, 0)
var userLeave chan User = make(chan User, 0)

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

func userHandler() {
	for {
		select {
		case user := <- userJoin:
			fmt.Println("Joined: " + user.connection.RemoteAddr().String())
			user.Start();
			users = append(users, user)
		case user := <- userLeave:
			fmt.Println("Left: " + user.connection.RemoteAddr().String())
			user.Close();//Maybe not the best idea to do this here, if there is a queue we might try read/write to the connection
			for p, u := range users {
				if(user == u) {
					users = append(users[:p], users[p + 1:]...)
				}
			}
		}
	}
}

func handleConnection(connection net.Conn) {
	user := User{connection, make(chan string, 0)}
	userJoin <- user
}