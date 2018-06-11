package main

import (
	"./server"
	"./client"
	"time"
)

func main() {
	server := gochatserver.New()
	server.Start()
	
	client2 := gochatclient.New()
	client2.Start("Bob")
	client2.SendMessage("Hello!")
	time.Sleep(time.Second * 1)
	client2.Stop()
	for {
	
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
