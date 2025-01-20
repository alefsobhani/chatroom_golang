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

1. **Build and Run with Docker Compose: Run the following command to build the Docker images and start the services:**

```bash
docker-compose up --build
```

2. **Access the Chat Server: The chat server will be available on port 8080. Clients can connect using:**

```bash
go run client.go
```

Test Commands

- Type messages in one client, and they appear in all connected clients.
- Use `/fusers` to list active users.

3. **Stop Services: To stop the containers, use:**

```bash
docker-compose down
```

## Verification

- NATS Management Interface: 

Open http://localhost:8222 in your browser to view the NATS server management interface.


- Client Interaction: 

Connect multiple clients and verify the chat functionality works as expected.
