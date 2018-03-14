package gochatclient

import (
	"net"
	"fmt"
	"strings"
)

type client struct {
	queue chan string
	connection net.Conn
}

func New() client {
	return client {
		make(chan string, 0),
		nil,
	}
}

func (client *client) Start() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if(err != nil) {
		panic(err)
	}
	client.connection = connection
	go client.handleConnectionWrite(connection)
	go client.handleConnectionRead(connection)
}

func (client *client) SendMessage(message string) {
	client.queue <- message;
}

func (client *client) handleConnectionWrite(connection net.Conn) {
	for {
		select {
		case str := <- client.queue:
			array := []byte(str[:len(str)])
			_, err := connection.Write(array)
			if(err != nil) {
				panic(err)
			}
		}	
	}
}

func (client *client) handleConnectionRead(connection net.Conn) {
	for {
		array := make([]byte, 1024);
		n, err := connection.Read(array)
		if(err != nil) {
			if(strings.HasSuffix(err.Error(), "An existing connection was forcibly closed by the remote host.")) {
				fmt.Println("Connection closed: " + connection.RemoteAddr().String())
				connection.Close();
				return
			}else {
				panic(err)
			}
		}
		message := string(array[:n]);
		fmt.Println(message);
	}
}