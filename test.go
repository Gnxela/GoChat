package main

import (
	"./server"
	"./client"
	"time"
)

func main() {
	server := gochatserver.New()
	client1 := gochatclient.New();
	client2 := gochatclient.New();
	server.Start();
	client1.Start("Alex");
	client2.Start("Bob");
	for {
		client1.SendMessage("Hello!")
		client2.SendMessage("Hello!")
		time.Sleep(time.Second * 1)
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