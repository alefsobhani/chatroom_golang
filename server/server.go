package main

import (
	"bufio"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
)

type ChatServer struct {
	users    map[string]net.Conn
	usersMu  sync.RWMutex
	natsConn *nats.Conn
	logger   *logrus.Logger
}

func NewChatServer() *ChatServer {
	// Initialize structured logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Read NATS URL from environment variable, default to localhost
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.Fatalf("Failed to connect to NATS: %v", err)
	}

	logger.Info("Chat server initialized")
	return &ChatServer{
		users:    make(map[string]net.Conn),
		natsConn: nc,
		logger:   logger,
	}
}

func (cs *ChatServer) broadcastMessage(message string) {
	cs.logger.WithFields(logrus.Fields{
		"type":    "broadcast",
		"message": message,
	}).Info("Broadcasting message")

	err := cs.natsConn.Publish("chat", []byte(message))
	if err != nil {
		cs.logger.WithError(err).Error("Failed to broadcast message")
	}
}

func (cs *ChatServer) listUsers(conn net.Conn) {
	cs.usersMu.RLock()
	defer cs.usersMu.RUnlock()

	cs.logger.WithField("client", conn.RemoteAddr().String()).Info("Listing active users")

	message := "Active users:\n"
	for addr := range cs.users {
		message += fmt.Sprintf(" - %s\n", addr)
	}

	_, err := conn.Write([]byte(message))
	if err != nil {
		cs.logger.WithError(err).WithField("client", conn.RemoteAddr().String()).Error("Failed to send user list")
	}
}

func (cs *ChatServer) handleClient(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	cs.logger.WithField("client", clientAddr).Info("New client connected")

	defer func() {
		cs.usersMu.Lock()
		delete(cs.users, clientAddr)
		cs.usersMu.Unlock()

		cs.broadcastMessage(fmt.Sprintf("%s left the chat\n", clientAddr))
		cs.logger.WithField("client", clientAddr).Info("Client disconnected")
		conn.Close()
	}()

	// Add the client to the active users
	cs.usersMu.Lock()
	cs.users[clientAddr] = conn
	cs.usersMu.Unlock()

	// Notify other users
	cs.broadcastMessage(fmt.Sprintf("%s joined the chat\n", clientAddr))

	// Subscribe to chat messages from NATS
	sub, err := cs.natsConn.SubscribeSync("chat")
	if err != nil {
		cs.logger.WithError(err).Error("Failed to subscribe to NATS")
		return
	}
	defer sub.Unsubscribe()

	// Goroutine to listen for NATS messages
	go func() {
		for {
			msg, err := sub.NextMsg(0)
			if err != nil {
				cs.logger.WithError(err).Error("Error reading NATS message")
				return
			}
			_, err = conn.Write(msg.Data)
			if err != nil {
				cs.logger.WithError(err).WithField("client", clientAddr).Error("Failed to send NATS message to client")
				return
			}
		}
	}()

	// Handle client input
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		cs.logger.WithFields(logrus.Fields{
			"client":  clientAddr,
			"message": text,
		}).Info("Received message")

		if text == "/fusers" {
			cs.listUsers(conn)
		} else {
			err := cs.natsConn.Publish("chat", []byte(fmt.Sprintf("%s: %s\n", clientAddr, text)))
			if err != nil {
				cs.logger.WithError(err).WithField("client", clientAddr).Error("Failed to publish message to NATS")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		cs.logger.WithError(err).WithField("client", clientAddr).Error("Error reading client input")
	}
}

func (cs *ChatServer) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		cs.logger.WithError(err).Fatal("Failed to start server")
	}
	defer listener.Close()

	cs.logger.WithField("port", port).Info("Server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			cs.logger.WithError(err).Error("Connection error")
			continue
		}
		go cs.handleClient(conn)
	}
}

func main() {
	server := NewChatServer()
	server.Start("8080")
}
