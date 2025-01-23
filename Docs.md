# Technical Documentation 

## **Overview**

This project implements a simple chat application using TCP sockets and NATS messaging. The system includes a **server** that handles client connections and broadcasts messages, and a **client** that allows users to send and receive messages in real time.

---

## **Server**

### **File**: `server.go`

### **Purpose**
The server is responsible for:
- Managing client connections.
- Broadcasting messages between clients using NATS messaging.
- Providing commands for specific functionalities, such as listing active users.

### **Features**
- **Broadcasting Messages**: Sends messages from one client to all connected clients.
- **User Commands**:
    - `/users`: Lists all currently connected clients.
- **Structured Logging**: Logs client activities, messages, and errors using the `logrus` library.
- **Graceful Client Handling**:
    - Detects client disconnections and removes them from the active user list.
    - Logs errors during client interaction for better traceability.

---

## **Client**

### **File**: `client.go`

### **Purpose**
The client provides an interface for users to:
- Connect to the chat server.
- Send and receive messages in real time.

### **Features**
- **Interactive Chat**: Allows users to type and send messages.
- **Server Commands**:
    - `/users`: Displays the list of connected users.
- **Error Handling**:
    - Detects and handles server disconnections gracefully.
    - Displays errors for invalid inputs or failed message deliveries.

---

## **Usage**

### **Running the Server**
1. Ensure that the **NATS server** is running on the default port (`4222`).
    - Start the NATS server:
      ```bash
      nats-server
      ```
2. Run the chat server:
   ```bash
   go run server.go
   ```
3. By default, the server listens on port `8080`. You can change the port by modifying the `Start` method:
   ```go
   server.Start("your_port")
   ```

### **Running the Client**
1. Connect to the server:
   ```bash
   go run client.go
   ```
2. Enter messages in the terminal to send them to the chat.
3. Use the `/users` command to view active users.

---

## **Commands**

| Command    | Description                                 |
|------------|---------------------------------------------|
| `/users`   | Lists all currently connected clients.      |
| `<message>`| Sends the entered message to all clients.   |

---

## **Configuration**

- **Environment Variables**:
    - `NATS_URL`: Specifies the URL for connecting to the NATS server. Defaults to `nats://localhost:4222` if not set.

---

## **Logging**
The server uses `logrus` for structured logging:
- **Info Logs**: Tracks client connections, disconnections, and sent messages.
- **Error Logs**: Captures failures, such as disconnections or message publishing errors.
- **Debug Logs**: Can be enabled to provide detailed traceability.

Example log output:
```
INFO[2025-01-23T15:45:12Z] New client connected                      client=192.168.1.5:12345
INFO[2025-01-23T15:45:15Z] Broadcasting message                      message="192.168.1.5:12345: Hello, World!"
INFO[2025-01-23T15:45:30Z] Client disconnected                       client=192.168.1.5:12345
```

---

## **Testing**

### **Server Tests**
Run the server tests:
```bash
go test -v server_test.go
```

### **Client Tests**
Ensure the client can:
- Connect to the server.
- Send and receive messages.
- Handle server disconnections.

---