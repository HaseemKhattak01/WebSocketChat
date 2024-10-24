package main

import (
	"fmt"
	"net/http"
	"websocket-chat/connections"
)

func main() {
	http.HandleFunc("/", connections.HomePage)
	http.HandleFunc("/ws", connections.HandleConnections)

	go connections.HandleMessages()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
