package gochatserver

import (
	"net"
	"strings"
	"../common"
)

type User struct {
	server *server
	connection net.Conn
	data []byte
	queue chan common.PacketMessage
	name string
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
				user.server.userLeave <- user;
				return
			}else {
				panic(err)
			}
		}
		user.data = append(user.data, array[:n]...);//Needs a read write lock, but will add later
		user.server.netManager.ReadData(&user.data)
	}
}

func (user User) handleConnectionWrite() {
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