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

- Server: A central server that connects users, listens for incoming messages, and broadcasts messages via NATS.
- Client: A CLI-based client that allows users to join the chat, send messages, and view the list of active users.

## Packages

- NATS: Use nats.go for NATS integration.
- CLI
- Concurrency: Use Goroutines for handling multiple clients concurrently.
- Networking

## Setup

1. **Run NATS server (using Docker):**
   ```bash
   docker-compose up -d
   