package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
)

type Client struct {
	username string
	conn     *websocket.Conn
	mail     chan []byte
	hub      *Hub
}

func (c *Client) ReadLoop() {
	for {
		_, msg, err := c.conn.ReadMessage()

		log.Println(string(msg))

		if err != nil {
			log.Println(err)
			return
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) WriteLoop() {
	for {
		writer, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println(err)
			return
		}
		msg, ok := <-c.mail
		if !ok {
			return
		}
		log.Println(string(msg))
		writer.Write(msg)
		if l := len(c.mail); l > 0 {
			for msg := range c.mail {
				writer.Write(msg)
			}
		}
		writer.Close()
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start upgrage")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//r.ParseForm()
	//username := r.Form.Get("username")

	client := &Client{
		conn:     conn,
	//	username: username,
		mail:     make(chan []byte, 1024),
		hub:      hub,
	}

	hub.register <- client
	go client.WriteLoop()
	go client.ReadLoop()
	fmt.Println("client online")
}
