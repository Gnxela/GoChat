package gochatclient

import (
	"net"
	"fmt"
	"strings"
	"../common"
	
	"github.com/Gnxela/GnPacket/GnPacket"
)

type client struct {
	queue chan common.PacketMessage
	connection net.Conn
	netManager GnPacket.NetManager
	data []byte
}

func New() client {
	return client {
		make(chan common.PacketMessage, 30),
		nil,
		GnPacket.New(100),
		nil,
	}
}

func (client *client) Start() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if(err != nil) {
		panic(err)
	}
	client.connection = connection
	
	client.netManager.AddHandler(1, client.handleMessage)
	
	go client.handleConnectionWrite()
	go client.handleConnectionRead()
}

func (client *client) handleMessage(packet GnPacket.GnPacket) bool {
	message := common.PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	fmt.Println(message.Message)
	return true
}

func (client *client) SendMessage(message string) {
	client.queue <- common.NewPacketMessage(message);
}

func (client *client) handleConnectionWrite() {
	for {
		select {
		case packet := <- client.queue:
			array := packet.Write(&packet)
			_, err := client.connection.Write(array)
			if(err != nil) {
				panic(err)
			}
		}	
	}
}

func (client *client) handleConnectionRead() {
	for {
		array := make([]byte, 1024);
		n, err := client.connection.Read(array)
		if(err != nil) {
			if(strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.")) {
				fmt.Println("Connection closed: " + client.connection.RemoteAddr().String())
				client.connection.Close();
				return
			}else {
				panic(err)
			}
		}
		client.data = append(client.data, array[:n]...)
		client.netManager.ReadData(&client.data)
	}
}