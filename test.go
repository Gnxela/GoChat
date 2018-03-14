package main

import (
	"./server"
	"./client"
	"time"
)

func main() {
	server := gochatserver.New()
	client := gochatclient.New();
	server.Start();
	client.Start();
	for {
		client.SendMessage("Hello!")
		time.Sleep(time.Second * 5)
	}
	/*
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if(err != nil) {
			panic(err)
		}
		client.SendMessage(str)
	}
	*/
}