package chat

import (
	"log"
)

type Hub struct {
	broadcast  chan []byte
	register   chan *Client
	unRegister chan *Client
	clients    map[*Client]bool
}

var hub = &Hub{
	broadcast:  make(chan []byte, 1024),
	register:   make(chan *Client, 1024),
	unRegister: make(chan *Client, 1024),
	clients:    make(map[*Client]bool),
}

func init() {
	go Run()
}

func Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unRegister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.mail)
			}
		case msg := <-hub.broadcast:
			log.Println(string(msg))
			for client := range hub.clients {
				select {
				case client.mail <- msg:
				default:
					close(client.mail)
					delete(hub.clients, client)
				}
			}
		}
	}
}
