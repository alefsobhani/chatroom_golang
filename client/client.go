package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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

func (cc *ChatClient) sendMessages() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		_, err = cc.conn.Write([]byte(text))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	}
}

func (cc *ChatClient) receiveMessages() {
	scanner := bufio.NewScanner(cc.conn)
	for scanner.Scan() {
		fmt.Print(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from server: %v", err)
	}
}
