package main

import (
	"log"
	"net"
)

type ChatClient struct {
	conn net.Conn
}

func NewChatClient(address string) *ChatClient {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	return &ChatClient{conn: conn}
}
