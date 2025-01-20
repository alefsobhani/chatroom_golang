# Chatroom Application

A simple command-line chatroom application written in Go using NATS for real-time messaging.

## Features
- Join the chatroom and send messages.
- View the list of active users using `/fusers`.
- Real-time messaging for all users in the chatroom.

## Requirements
- Go 1.18+
- Docker (for NATS server)

## Architecture

- The Server listens for incoming connections, handles users, and broadcasts messages using NATS.
- The client application will connect to the server over TCP and handle input/output using the CLI.

## Packages

- NATS: Use nats.go for NATS integration.
- CLI
- Concurrency: Use Goroutines for handling multiple clients concurrently.
- Networking

## Setup

1. **Run each command on separate terminal:**
   ```bash
   docker-compose up -d
   go run ./server/server.go
   go run ./server/client.go

## Test Commands

- Type messages in one client, and they appear in all connected clients.
- Use `/fusers` to list active users.