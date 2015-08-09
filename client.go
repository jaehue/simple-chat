package main

import (
	"log"
	"time"
)

const messageBufferSize = 256

type Client struct {
	Messager
	Name      string
	AvatarURL string
	send      chan *Message
}

func newClient(messager Messager, name, avatarUrl string) *Client {
	c := &Client{
		Messager:  messager,
		Name:      name,
		AvatarURL: avatarUrl,
		send:      make(chan *Message, messageBufferSize),
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
	var msg Message
	if err := c.ReadChatMessage(&msg); err != nil {
		return nil, err
	}
	msg.CreatedAt = time.Now()
	msg.Name = c.Name
	msg.AvatarURL = c.AvatarURL
	log.Println("read from websocket:", msg)
	return &msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)
	return c.WriteChatMessage(m)
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
				c.WriteChatMessage(&Message{MessageType: ChatError, Content: err.Error()})
			}
		}

	}
	log.Printf("close connection. Client: %v", c)
	c.Close()
}

func (c *Client) writeLoop() {
	for msg := range c.send {
		c.write(msg)
	}
	log.Printf("close connection. Client: %v", c)
	c.Close()
}

func (c *Client) doRoomAction(m *Message) error {
	if m.Content == "enter" {
		if r := FindRoom(m.RoomId); r != nil {
			c.join(r)
		}
	} else if m.Content == "leave" {
		if r := FindRoom(m.RoomId); r != nil {
			c.leave(r)
		}
	}
	return nil
}
