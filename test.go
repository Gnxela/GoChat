package main

import (
	"./server"
)

func main() {
	server := gochatserver.New()
	server.Start();
}