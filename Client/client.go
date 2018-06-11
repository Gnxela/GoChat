package gochatclient

import (
	"net"
	"fmt"
	"strings"
	"../common"
	
	"github.com/Gnxela/GnPacket/GnPacket"
)

type client struct {
	queue			chan common.PacketMessage
	connection		net.Conn
	netManager		GnPacket.NetManager
	data			[]byte
	Name			string
	stopChannel		chan int8
}

func New() client {
	return client {
		make(chan common.PacketMessage, 30),
		nil,
		GnPacket.New(100),
		nil,
		"",
		make(chan int8, 10),//Reserving 10 spaces, really should not need this but meh
	}
}

func (client *client) Start(name string) {
	connection, err := net.Dial("tcp", "localhost:8080")
	if(err != nil) {
		panic(err)
	}
	client.connection = connection
	
	client.netManager.AddHandler(0, client.handshakeHandler)
	
	client.Name = name
	client.sendHandshake()
	
	go client.handleConnectionWrite()
	go client.handleConnectionRead()
}

func (client *client) Stop() {
	client.stopChannel <- 0
	client.connection.Close()
}

func (client *client) sendHandshake() {
	packet := common.NewPacketHandshake(client.Name)
	array := packet.Write(&packet)
	client.connection.Write(array)
}

func (client *client) handshakeHandler(packet GnPacket.GnPacket) bool {
	handshake := common.PacketHandshake{&packet, ""}
	handshake.Deserialize(packet.Data)
	
	if (handshake.Name == "") {
		client.netManager.RemoveHandler(0, client.handshakeHandler)
		client.netManager.AddHandler(1, client.handleMessage)
		client.recycleUnhandledPackets()
	} else {
		panic("Server rejected handshake: " + handshake.Name);
	}
	
	return false//No other handlers should ever really get this.
}

func (client *client) recycleUnhandledPackets() {
	num := len(client.netManager.UnhandledQueue)
	L: for i := 0; i < num; i++ {
		select {
		case packet := <- client.netManager.UnhandledQueue:
			client.netManager.DispatchPacket(packet)
		default:
			break L//Breaks the for loop.
		}
	}
}

func (client *client) handleMessage(packet GnPacket.GnPacket) bool {
	message := common.PacketMessage{&packet, ""}
	message.Deserialize(packet.Data)
	fmt.Println(client.Name + "|" + message.Message)
	return true
}

func (client *client) SendMessage(message string) {
	client.queue <- common.NewPacketMessage(message);
}

func (client *client) handleConnectionWrite() {
	L: for {
		select {
		case <- client.stopChannel:
			client.stopChannel <- 0
			break L
		case packet := <- client.queue:
			array := packet.Write(&packet)
			_, err := client.connection.Write(array)
			if(err != nil) {
				panic(err)
			}
			break
		}
	}
}

func (client *client) handleConnectionRead() {
	L: for {
		select {
		case <- client.stopChannel:
			client.stopChannel <- 0
			break L
		default:
			array := make([]byte, 1024);
			n, err := client.connection.Read(array)
			if(err != nil) {
				if strings.HasSuffix(err.Error(), "use of closed network connection") {//Honestly spent 20 minutes looking for a better way to do this. The error is from "poll", ErrNetClosing. Can't import poll though
					client.Stop()
				} else if strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.") {
					client.Stop()
					return
				} else {
					fmt.Println(err)
				}
			}
			client.data = append(client.data, array[:n]...)
			client.netManager.ReadData(&client.data)
			break
		}
	}
}