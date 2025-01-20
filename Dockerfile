FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o chat-server server.go
EXPOSE 8080
CMD ["./chat-server"]
