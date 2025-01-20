package main

import (
	"fmt"
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

func (cs *ChatServer) broadcastMessage(message string) {
	err := cs.natsConn.Publish("chat", []byte(message))
	if err != nil {
		log.Fatalf("Failed to broadcast message to chatserver: %v", err)
	}
}

func (cs *ChatServer) listUsers(conn net.Conn) {
	cs.usersMu.RLock()
	defer cs.usersMu.RUnlock()

	_, err := conn.Write([]byte("Active users:\n"))
	for addr := range cs.users {
		_, err = conn.Write([]byte(addr + "\n"))
	}

	if err != nil {
		log.Fatalf("Failed to make list users from chatserver: %v", err)
	}
}

// TODO: handle errors of this function in future
func (cs *ChatServer) handleClient(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	cs.usersMu.Lock()
	cs.users[clientAddr] = conn
	cs.usersMu.Unlock()

	cs.broadcastMessage(fmt.Sprintf("%s joined the chat\n", clientAddr))
	
	sub, err := cs.natsConn.SubscribeSync("chat")
	if err != nil {
		log.Printf("Failed to subscribe to NATS: %v", err)
		return
	}
	defer sub.Unsubscribe()

}

func main() {

}
