package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
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

	var userData map[string]interface{}
	if authCookie, err := r.Cookie("auth"); err == nil {
		userData = objx.MustFromBase64(authCookie.Value)
	}

	newClient(&messager, userData["name"].(string), userData["avatar_url"].(string))
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
