package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const socketBufferSize = 1024

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	log.Println("connect socket. addr: ", socket.RemoteAddr())
	newClient(socket)
}
