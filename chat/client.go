package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}

	err := c.socket.Close()
	if err != nil {
		log.Fatal("Failure socket close")
	}
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}

	err := c.socket.Close()
	if err != nil {
		log.Fatal("Failure socket close")
	}
}