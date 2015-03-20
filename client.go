package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const messageBufferSize = 256

type Client struct {
	conn *websocket.Conn
	send chan *Message
}

func newClient(conn *websocket.Conn) *Client {
	c := &Client{
		conn: conn,
		send: make(chan *Message, messageBufferSize),
	}
	c.listen()
	return c
}

func (c *Client) listen() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) join(r *Room) {
	r.join <- c
}

func (c *Client) leave(r *Room) {
	r.leave <- c
}

func (c *Client) read() (*Message, error) {
	var msg *Message
	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}
	msg.CreatedAt = time.Now()
	log.Println("read from websocket:", msg)
	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)
	return c.conn.WriteJSON(m)
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}
		if m.MessageType == TextMessage {
			if r := FindRoom(m.RoomId); r != nil {
				r.forward <- m
			}
			continue
		} else if m.MessageType == RoomAction {
			if err := c.doRoomAction(m); err != nil {
				log.Println(err)
				c.conn.WriteJSON(struct{ Status, Message string }{Status: "Fail", Message: err.Error()})
			}
		}

	}
	log.Printf("close connection. addr: %s", c.conn.RemoteAddr())
	c.conn.Close()
}

func (c *Client) writeLoop() {
	for msg := range c.send {
		c.write(msg)
	}
	log.Printf("close connection. Client: %v", c.conn.RemoteAddr())
	c.conn.Close()
}

func (c *Client) doRoomAction(m *Message) error {
	enterRgx := regexp.MustCompile("enter room:.+")
	leaveRgx := regexp.MustCompile("leave room:.+")

	if str := enterRgx.FindString(m.Content); len(str) > 0 {
		id := strings.TrimSpace(str[len("enter room:"):])
		if len(id) > 0 {
			i, err := strconv.Atoi(id)
			if err != nil {
				c.conn.WriteMessage(websocket.TextMessage, []byte("Room id is invalid"))
				return fmt.Errorf("Room id is invalid. error: %v", err)
			}
			if r := FindRoom(i); r != nil {
				c.join(r)
			}
		}
	}

	if str := leaveRgx.FindString(m.Content); len(str) > 0 {
		id := strings.TrimSpace(str[len("leave room:"):])
		if len(id) > 0 {
			i, err := strconv.Atoi(id)
			if err != nil {
				c.conn.WriteMessage(websocket.TextMessage, []byte("Room id is invalid"))
				return fmt.Errorf("Room id is invalid. error: %v", err)
			}
			if r := FindRoom(i); r != nil {
				c.leave(r)
			}
		}
	}
	return nil
}
