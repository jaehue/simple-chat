package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const socketBufferSize = 1024

type webSocketConn websocket.Conn

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
	var messager webSocketConn = webSocketConn(*socket)
	newClient(&messager)
}

func (c *webSocketConn) ReadChatMessage(msg *Message) error {
	var conn websocket.Conn = websocket.Conn(*c)

	if err := conn.ReadJSON(&msg); err != nil {
		return err
	}
	msg.CreatedAt = time.Now()
	log.Println("read from websocket:", msg)
	return nil
}

func (c *webSocketConn) WriteChatMessage(msg *Message) error {
	var conn websocket.Conn = websocket.Conn(*c)
	return conn.WriteJSON(msg)
}

func (c *webSocketConn) Close() error {
	var conn websocket.Conn = websocket.Conn(*c)
	return conn.Close()
}

func (c *webSocketConn) String() string {
	var conn websocket.Conn = websocket.Conn(*c)
	return conn.RemoteAddr().String()
}
