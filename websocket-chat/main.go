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

	fmt.Println("Server started on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
