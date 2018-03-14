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
		make(chan string, 30),
		nil,
	}
}

func (client *client) Start() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if(err != nil) {
		panic(err)
	}
	client.connection = connection
	go client.handleConnectionWrite()
	go client.handleConnectionRead()
}

func (client *client) SendMessage(message string) {
	client.queue <- message;
}

func (client *client) handleConnectionWrite() {
	for {
		select {
		case str := <- client.queue:
			array := []byte(str[:len(str)])
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
		message := string(array[:n]);
		fmt.Println(message);
	}
}