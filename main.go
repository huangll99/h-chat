package main

import (
	"net/http"
	"h-chat/chat"
	"log"
)

func main() {
	http.HandleFunc("/ws", chat.ServeWs)

	log.Println("chat server started.")
	http.ListenAndServe(":8888", nil)
}
