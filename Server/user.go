package gochatserver

import (
	"net"
	"strings"
)

type User struct {
	server server
	connection net.Conn
	queue chan string
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
		str := strings.TrimSpace(string(array[:n]))
		message := Message{&user, str}
		user.server.userMessage <- message
	}
}

func (user User) handleConnectionWrite() {
	for {
		select {
		case str := <- user.queue:
			array := []byte(str[:len(str)])
			_, err := user.connection.Write(array)
			if(err != nil) {//
				if(!strings.HasSuffix(err.Error(), "use of closed network connection")) {
					panic(err)
				}
			}
		}	
	}
}