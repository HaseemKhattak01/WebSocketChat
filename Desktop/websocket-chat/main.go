package main

import (
	"fmt"
	"net/http"
	"websocket-chat/connections"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Initialize Redis client
	connections.InitRedis()

	http.HandleFunc("/", connections.HomePage)
	http.HandleFunc("/ws", connections.HandleConnections)
	http.HandleFunc("/health", healthCheck)

	go connections.HandleMessages()

	fmt.Println("Server started on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
