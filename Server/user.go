package gochatserver

import (
	"net"
	"strings"
	"bytes"
	"../common"
	
	"github.com/Gnxela/GnPacket/GnPacket"
)

type User struct {
	server *server
	connection net.Conn
	netManager GnPacket.NetManager
	data []byte
	queue chan common.PacketMessage
	Name string
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
	
	user.Name = handshake.Name;
	
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
	user.recycleUnhandledPackets()
	user.server.AddUser(user)
	
	return false
}

func (user *User) recycleUnhandledPackets() {
	num := len(user.netManager.UnhandledQueue)
	L: for i := 0; i < num; i++ {
		select {
		case packet := <- user.netManager.UnhandledQueue:
			user.netManager.DispatchPacket(packet)
		default:
			break L//Breaks the for loop.
		}
	}
}

func (user *User) handleMessage(packet GnPacket.GnPacket) bool {
	message := common.PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	var buffer bytes.Buffer
	buffer.WriteString("<")
	buffer.WriteString(user.Name)
	buffer.WriteString("> ")
	buffer.WriteString(message.Message)
	user.server.SendMessage(buffer.String())
	return true
}


func (user *User) handleConnectionRead() {
	for {
		array := make([]byte, 1024);
		n, err := user.connection.Read(array)
		if (err != nil) {
			if(strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.")) {//TODO rather than string comparison there's probably an error object I can compare to.
				user.server.RemoveUser(user);
				return
			}else {
				panic(err)
			}
		}
		user.data = append(user.data, array[:n]...);//TODO Needs a read write lock, but will add later
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