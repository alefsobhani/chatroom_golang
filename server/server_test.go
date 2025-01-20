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

func TestUserJoinAndLeave(t *testing.T) {
	server := startTestChatServer(t)
	defer server.natsConn.Close()

	client1 := connectTestClient(t)
	defer client1.Close()

	client2 := connectTestClient(t)
	defer client2.Close()

	// Read join messages from client1
	client1Scanner := bufio.NewScanner(client1)
	var messages []string
	var mu sync.Mutex
	go func() {
		for client1Scanner.Scan() {
			mu.Lock()
			messages = append(messages, client1Scanner.Text())
			mu.Unlock()
		}
	}()

	// Wait for client2's join message
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(messages) == 0 || !strings.Contains(messages[len(messages)-1], "joined the chat") {
		t.Errorf("Expected join message not received")
	}

	// Close client2 and verify leave message
	client2.Close()
	time.Sleep(500 * time.Millisecond)

	if len(messages) == 0 || !strings.Contains(messages[len(messages)-1], "left the chat") {
		t.Errorf("Expected leave message not received")
	}
}

func TestListActiveUsers(t *testing.T) {
	server := startTestChatServer(t)
	defer server.natsConn.Close()

	client1 := connectTestClient(t)
	defer client1.Close()

	client2 := connectTestClient(t)
	defer client2.Close()

	// Client1 requests the list of active users
	client1.Write([]byte("/fusers\n"))

	// Capture the response
	client1Scanner := bufio.NewScanner(client1)
	var userList []string
	for client1Scanner.Scan() {
		line := client1Scanner.Text()
		if strings.Contains(line, "Active users") || strings.Contains(line, "127.0.0.1") {
			userList = append(userList, line)
		}
		if len(userList) >= 2 { // We expect at least two active users
			break
		}
	}

	if len(userList) < 2 {
		t.Errorf("Expected at least 2 users in the active users list, got: %v", userList)
	}

	if !strings.Contains(userList[1], "127.0.0.1") {
		t.Errorf("User list does not contain expected client addresses")
	}
}
