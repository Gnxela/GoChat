package gochatserver

import (
	"net"
	"fmt"
	"bytes"
	"../common"
	
	"github.com/Gnxela/GnPacket/GnPacket"
)

type server struct {
	netManager GnPacket.NetManager
	users []User
	userJoin chan User 
	userLeave chan User
	userMessage chan Message
}

type Message struct {
	sender *User
	message string
}

func New() server {
	return server{
		GnPacket.New(100),
		make([]User, 0),
		make(chan User, 10),
		make(chan User, 10),
		make(chan Message, 100),
	}
}

func (server *server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if(err != nil) {
		panic(err)
	}
	
	server.netManager.AddHandler(1, server.handleMessage)
	
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
	user := User{server, connection, nil, make(chan common.PacketMessage, 30), "User"}
	server.userJoin <- user//Won't remove this yet, it's actually a nice way to ensure no concurrent errors
}

func (server *server) handleMessage(packet GnPacket.GnPacket) bool {
	message := common.PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	var buffer bytes.Buffer
	buffer.WriteString("<")
	//buffer.WriteString(message.sender.name)
	buffer.WriteString("TEST")
	buffer.WriteString("> ")
	buffer.WriteString(message.Message)
	server.SendMessage(buffer.String())//Handle all messages in a single routine so that we ensure that they are ordered correctly for all clients. "correctly" not nessesarily being the right order, but a consistant order
	return true
}

func (server *server) userHandler() {
	for {
		select {
		case user := <- server.userJoin:
			user.Start();
			server.users = append(server.users, user)
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