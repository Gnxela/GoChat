package gochatserver

import (
	"net"
	"fmt"
	"sync"
	"../common"

	"github.com/Gnxela/GnPacket/GnPacket"	
)

type server struct {
	users []*User
	usersLock *sync.Mutex
	sendMessageLock *sync.Mutex
}

type Message struct {
	sender *User
	message string
}

func New() server {
	return server{
		make([]*User, 0),
		&sync.Mutex{},
		&sync.Mutex{},
	}
}

func (server *server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if(err != nil) {
		panic(err)
	}
		
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
	user.Start();
}

func (server *server) AddUser(user *User) {
	server.usersLock.Lock()
	server.users = append(server.users, user)
	server.usersLock.Unlock()
	server.SendMessage(user.Name + " joined the server.")
}

func (server *server) RemoveUser(user *User) {
	server.usersLock.Lock()
	for p, u := range server.users {
		if(user == u) {
			server.users = append(server.users[:p], server.users[p + 1:]...)
		}
	}
	server.usersLock.Unlock()

	user.Close();//Maybe not the best idea to do this here, if there is a queue we might try read/write to the connection. I am correct, there is a rare error when disconnecting.
	server.SendMessage(user.Name + " left the server.")
}

func (server *server) SendMessage(str string) {
	server.sendMessageLock.Lock()
	fmt.Println("> " + str)
	for _, u := range server.users {
		u.queue <- common.NewPacketMessage(str)
	}
	server.sendMessageLock.Unlock()
}