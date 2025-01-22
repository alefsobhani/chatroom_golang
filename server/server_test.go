package main

import (
	"bufio"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func setupTestServer(t *testing.T) (*ChatServer, string) {
	// Create a mock NATS server URL for testing
	natsURL := "nats://127.0.0.1:4222"

	// Set the environment variable for the NATS URL
	os.Setenv("NATS_URL", natsURL)

	// Start the mock NATS server
	nc, err := nats.Connect(natsURL)
	assert.NoError(t, err)

	// Create a ChatServer instance
	server := &ChatServer{
		users:    make(map[string]net.Conn),
		usersMu:  sync.RWMutex{},
		natsConn: nc,
		logger:   NewChatServer().logger,
	}

	return server, natsURL
}

func TestServerStart(t *testing.T) {
	server, _ := setupTestServer(t)

	// Start the server in a separate goroutine
	go func() {
		server.Start("9090")
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	conn.Close()
}

func TestBroadcastMessage(t *testing.T) {
	server, _ := setupTestServer(t)

	// Capture the broadcast messages
	var receivedMessages []string
	sub, err := server.natsConn.Subscribe("chat", func(msg *nats.Msg) {
		receivedMessages = append(receivedMessages, string(msg.Data))
	})
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	// Broadcast a message
	server.broadcastMessage("Test message")

	// Allow some time for the message to be received
	time.Sleep(500 * time.Millisecond)

	// Validate the broadcast message
	assert.Len(t, receivedMessages, 1)
	assert.Equal(t, "Test message", receivedMessages[0])
}

func TestListUsers(t *testing.T) {
	server, _ := setupTestServer(t)

	// Start the server in a separate goroutine
	go func() {
		server.Start("9091")
	}()
	time.Sleep(1 * time.Second)

	// Connect multiple clients
	client1, err := net.Dial("tcp", "127.0.0.1:9091")
	assert.NoError(t, err)
	defer client1.Close()

	client2, err := net.Dial("tcp", "127.0.0.1:9091")
	assert.NoError(t, err)
	defer client2.Close()

	// Send the `/fusers` command from client1
	fmt.Fprintln(client1, "/fusers")

	// Read the response
	scanner := bufio.NewScanner(client1)
	scanner.Scan()
	output := scanner.Text()

	// Validate the list of active users
	assert.Contains(t, output, "Active users")
	assert.Contains(t, output, client1.RemoteAddr().String())
	assert.Contains(t, output, client2.RemoteAddr().String())
}

func TestHandleClientDisconnection(t *testing.T) {
	server, _ := setupTestServer(t)

	// Start the server in a separate goroutine
	go func() {
		server.Start("9092")
	}()
	time.Sleep(1 * time.Second)

	// Connect a client
	client, err := net.Dial("tcp", "127.0.0.1:9092")
	assert.NoError(t, err)

	clientAddr := client.RemoteAddr().String()

	// Ensure the client is added to the users list
	server.usersMu.RLock()
	_, exists := server.users[clientAddr]
	server.usersMu.RUnlock()
	assert.True(t, exists)

	// Disconnect the client
	client.Close()
	time.Sleep(500 * time.Millisecond)

	// Ensure the client is removed from the users list
	server.usersMu.RLock()
	_, exists = server.users[clientAddr]
	server.usersMu.RUnlock()
	assert.False(t, exists)
}

func TestClientMessageHandling(t *testing.T) {
	server, _ := setupTestServer(t)

	// Start the server in a separate goroutine
	go func() {
		server.Start("9093")
	}()
	time.Sleep(1 * time.Second)

	// Connect a client
	client, err := net.Dial("tcp", "127.0.0.1:9093")
	assert.NoError(t, err)
	defer client.Close()

	// Mock NATS subscriber to capture messages
	var receivedMessages []string
	sub, err := server.natsConn.Subscribe("chat", func(msg *nats.Msg) {
		receivedMessages = append(receivedMessages, string(msg.Data))
	})
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	// Send a message from the client
	testMessage := "Hello, ChatServer!"
	fmt.Fprintln(client, testMessage)

	// Allow some time for the message to propagate
	time.Sleep(500 * time.Millisecond)

	// Validate the received message
	assert.Len(t, receivedMessages, 1)
	assert.Contains(t, receivedMessages[0], testMessage)
}
