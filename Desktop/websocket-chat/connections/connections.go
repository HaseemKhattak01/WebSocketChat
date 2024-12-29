package connections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"websocket-chat/models"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	Clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	RedisClient *redis.Client
	Broadcast   = make(chan models.Message)
	ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "redis-service:6379",
	})
	go subscribeToRedis()
}

func subscribeToRedis() {
	sub := RedisClient.Subscribe(ctx, "chat_channel")
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Printf("Error receiving message from Redis: %v\n", err)
			continue
		}

		var message models.Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		for client := range Clients {
			if err := client.WriteJSON(message); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	Clients[conn] = true
	defer delete(Clients, conn)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading JSON: %v\n", err)
			break
		}

		msgJSON, _ := json.Marshal(msg)
		RedisClient.Publish(ctx, "chat_channel", msgJSON)
	}
}

func HandleMessages() {
	for {
		msg := <-Broadcast

		msgJSON, _ := json.Marshal(msg)
		RedisClient.Publish(ctx, "chat_channel", msgJSON)
	}
}
