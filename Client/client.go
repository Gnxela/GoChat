package main

import (
	"net"
	"fmt"
	"sync"
	"bufio"
	"os"
	"strings"
)

var wg sync.WaitGroup
var queue chan string

func main() {
	queue = make(chan string, 0)
	connection, err := net.Dial("tcp", "localhost:8080")
	if(err != nil) {
		panic(err)
	}
	fmt.Println("Client connected: " + connection.RemoteAddr().String())
	go handleConnectionWrite(connection)
	go handleConnectionRead(connection)
	
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if(err != nil) {
			panic(err)
		}
		queue <- str
	}
}

func handleConnectionWrite(connection net.Conn) {
	for {
		select {
		case str := <- queue:
			array := []byte(str[:len(str) - 1])
			_, err := connection.Write(array)
			if(err != nil) {
				panic(err)
			}
		}	
	}
}

func handleConnectionRead(connection net.Conn) {
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