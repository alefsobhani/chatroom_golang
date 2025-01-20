package main

import (
	"net"
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
