package gochatserver

import (
	"net"
	"strings"
	"bytes"
	"../common"
	
	"github.com/Gnxela/GnPacket/GnPacket"

	"fmt"
)

type User struct {
	server *server
	connection net.Conn
	netManager GnPacket.NetManager
	data []byte
	queue chan common.PacketMessage
	name string
}

func (user *User) Start() {
	user.netManager.AddHandler(0, user.handleHandshake)
	
	go user.handleConnectionRead()
	go user.handleConnectionWrite()
}

func (user *User) Close() {
	user.connection.Close()
}

//Returns true when user is accepted, false when a handshake is rejected.
func (user *User) handleHandshake(packet GnPacket.GnPacket) (bool) {
	//Read handshake data
	handshake := common.PacketHandshake{&packet, ""}
	handshake.Deserialize(packet.Data)
	
	//Validate data
	var userError string = ""
	
	user.name = handshake.Name;
	
	if handshake.Name == "" {
		userError = "Username invalid"
	}

	//Return handshake
	returnPacket := common.NewPacketHandshake(userError)
	returnArray := returnPacket.Write(&returnPacket)
	_, err := user.connection.Write(returnArray)
	if err != nil {
		panic(err)
	}
	
	//Setup user
	user.netManager.RemoveHandler(0, user.handleHandshake)
	user.netManager.AddHandler(1, user.handleMessage)
	fmt.Printf("Unhandled Packets After Handshake(S): %d\n", len(user.netManager.UnhandledQueue))
	
	return false
}

func (user *User) handleMessage(packet GnPacket.GnPacket) bool {
	message := common.PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	var buffer bytes.Buffer
	buffer.WriteString("<")
	//buffer.WriteString(message.sender.name)
	buffer.WriteString(user.name)
	buffer.WriteString("> ")
	buffer.WriteString(message.Message)
	user.server.SendMessage(buffer.String())//Handle all messages in a single routine so that we ensure that they are ordered correctly for all clients. "correctly" not nessesarily being the right order, but a consistant order
	return true
}


func (user *User) handleConnectionRead() {
	for {
		array := make([]byte, 1024);
		n, err := user.connection.Read(array)
		if (err != nil) {
			if(strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.")) {
				user.server.userLeave <- user;
				return
			}else {
				panic(err)
			}
		}
		user.data = append(user.data, array[:n]...);//Needs a read write lock, but will add later
		user.netManager.ReadData(&user.data)
	}
}

func (user *User) handleConnectionWrite() {
	for {
		select {
		case packet := <- user.queue:
			array := packet.Write(&packet)
			_, err := user.connection.Write(array)
			if(err != nil) {
				if(!strings.HasSuffix(err.Error(), "use of closed network connection")) {
					panic(err)
				}
			}
		}	
	}
}