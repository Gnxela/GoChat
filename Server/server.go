package gochatserver

import (
	"net"
	"fmt"
	"../common"

	"github.com/Gnxela/GnPacket/GnPacket"	
)

type server struct {
	users []*User
	userJoin chan *User 
	userLeave chan *User
	userMessage chan Message
}

type Message struct {
	sender *User
	message string
}

func New() server {
	return server{
		make([]*User, 0),
		make(chan *User, 10),
		make(chan *User, 10),
		make(chan Message, 100),
	}
}

func (server *server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if(err != nil) {
		panic(err)
	}
		
	go server.userHandler()
	go server.connectionListener(listener)
}

func (server *server) connectionListener(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if(err != nil) {
			panic(err)
		}
		go server.handleConnection(connection)
	}
}

func (server *server) handleConnection(connection net.Conn) {
	user := User{server, connection, GnPacket.New(100), nil, make(chan common.PacketMessage, 30), "User"}
	server.userJoin <- &user//TODO. Create a server add/remove user function or something. Won't remove this yet, it's actually a nice way to ensure no concurrent errors
}

func (server *server) userHandler() {
	for {
		select {
		case user := <- server.userJoin:
			user.Start();
			server.users = append(server.users, user)//TODO need to change what happens here. Should wait for handshake
			server.SendMessage("A user joined the server.")
		case user := <- server.userLeave:
			/*for p, u := range server.users {
				if(user == u) {
					server.users = append(server.users[:p], server.users[p + 1:]...)
				}
			}*/
			user.Close();//Maybe not the best idea to do this here, if there is a queue we might try read/write to the connection. I am correct, there is a rare error when disconnecting.
			server.SendMessage(user.name + " left the server.")
		}
	}
}

func (server *server) SendMessage(str string) {
	fmt.Println("> " + str);
	for _, u := range server.users {
		u.queue <- common.NewPacketMessage(str)
	}
}