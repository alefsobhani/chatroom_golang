package main

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func startTestChatServer(t *testing.T) *ChatServer {
	server := NewChatServer()
	go func() {
		server.Start("9090") // Use a different port for testing
	}()
	// Give the server some time to start
	time.Sleep(500 * time.Millisecond)
	return server
}

func connectTestClient(t *testing.T) net.Conn {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Fatalf("Failed to connect to the server: %v", err)
	}
	return conn
}

func TestBroadcastMessage(t *testing.T) {
	server := startTestChatServer(t)
	defer server.natsConn.Close()

	client1 := connectTestClient(t)
	defer client1.Close()

	client2 := connectTestClient(t)
	defer client2.Close()

	// Goroutine to listen for messages on client2
	var receivedMessages []string
	var mu sync.Mutex
	go func() {
		client2Scanner := bufio.NewScanner(client2)
		for client2Scanner.Scan() {
			mu.Lock()
			receivedMessages = append(receivedMessages, client2Scanner.Text())
			mu.Unlock()
		}
	}()

	// Client1 sends a message
	message := "Hello, World!"
	client1.Write([]byte(message + "\n"))

	// Wait for the message to propagate
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(receivedMessages) == 0 || !strings.Contains(receivedMessages[len(receivedMessages)-1], message) {
		t.Errorf("Broadcast message not received by other clients")
	}
}
