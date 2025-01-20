package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"net"
	"sync"
)

type ChatServer struct {
	users    map[string]net.Conn
	usersMu  sync.RWMutex
	natsConn *nats.Conn
}

func NewChatServer() *ChatServer {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	return &ChatServer{
		users:    make(map[string]net.Conn),
		natsConn: nc,
	}
}

func main() {
}
