package main

import (
	"./server"
	"time"
	"fmt"
)

func main() {
	server := gochatserver.New()
	go server.Start();
	for {
		server.SendMessage("Hello!")
		fmt.Println("Hello!")
		time.Sleep(time.Second * 5)
	}
}