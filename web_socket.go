package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/objx"
)

const socketBufferSize = 1024

type webSocketConn struct {
	websocket.Conn
	quit chan struct{}
}

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
	messager := webSocketConn{*socket, make(chan struct{})}

	var userData map[string]interface{}
	if authCookie, err := r.Cookie("auth"); err == nil {
		userData = objx.MustFromBase64(authCookie.Value)
	}

	newClient(&messager, userData["name"].(string), userData["avatar_url"].(string))
}

func (c *webSocketConn) ReadChatMessage(msg *Message) error {
	if err := c.Conn.ReadJSON(&msg); err != nil {
		return err
	}
	msg.CreatedAt = time.Now()
	log.Println("read from websocket:", msg)
	return nil
}

func (c *webSocketConn) WriteChatMessage(msg *Message) error {
	return c.Conn.WriteJSON(msg)
}

func (c *webSocketConn) Close() error {
	c.quit <- struct{}{}
	return c.Conn.Close()
}

func (c *webSocketConn) RemoteAddr() string {
	return c.Conn.RemoteAddr().String()
}

func (c *webSocketConn) Quit() <-chan struct{} {
	return c.quit
}
