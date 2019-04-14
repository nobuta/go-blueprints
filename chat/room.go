package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	forward chan []byte
	join chan *client
	leave chan *client
	clients map[*client]bool
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) run(){
	for {
		select {
		case client := <- r.join:

			log.Println("join")
			r.clients[client] = true
		case client := <- r.leave:

			delete(r.clients, client)
			close(client.send)
		case msg := <- r.forward:

			for client := range r.clients {
				select {
				case client.send <- msg:

					// 送信処理
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	log.Println(req.RemoteAddr, req.Header.Get("user-agent"))

	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client

	// leaveチャンネルにclient渡す(client.read()を抜けたタイミング(正常, 異常終了問わず)
	defer func() {
		r.leave <- client
	}()
	// websocketコネクションへメッセージ書き込みするgoroutineを作成
	go client.write()
	// websocketコネクションのメッセージ受信待ち続ける(ブロック)
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map [*client]bool),
	}
}