# Explanation of Tests

- TestUserJoinAndLeave: 

Starts the chat server and connects two clients.
Verifies that the join and leave notifications are broadcast correctly to other users.


- TestBroadcastMessage: 

Starts the chat server and connects two clients.
Verifies that a message sent by one client is received by another client.


- TestListActiveUsers:

Starts the chat server and connects two clients.
Verifies that the /fusers command lists all active users.


## How to Run Tests
Run the tests using the go test command:

1. Install dependencies:

```bash
go get -u github.com/nats-io/nats.go github.com/stretchr/testify
```

2. Run the test suite:

```bash
go test -v
```

The -v flag enables verbose output, so you can see the details of each test case.

- Concurrency Handling: 

Tests use sync.Mutex to handle shared data safely.


- Port Isolation: 

The tests run on a separate port (9090) to avoid conflicts with a production server.


- Reusable Functions: 

Functions like startTestChatServer and connectTestClient make the test setup reusable and clean.


**This test suite ensures all key features of the chat server work as expected. Let me know if you need additional tests or refinements!**