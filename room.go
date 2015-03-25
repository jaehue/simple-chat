package main

import (
	"log"
	"sync"
	"sync/atomic"
)

type Room struct {
	Id      int
	Name    string
	forward chan *Message
	join    chan *Client
	leave   chan *Client
	clients map[*Client]bool
}

var (
	mu        sync.Mutex
	maxRoomId int32 = 0
	rooms           = make([]*Room, 0)
)

func FindRoom(id int) *Room {
	for _, r := range rooms {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func NewRoom(name string) *Room {
	atomic.AddInt32(&maxRoomId, 1)
	r := &Room{
		Id:      int(maxRoomId),
		Name:    name,
		forward: make(chan *Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}

	mu.Lock()
	defer mu.Unlock()
	rooms = append(rooms, r)

	go r.run()
	return r
}

func (r *Room) broadcast(m *Message) {
	for client := range r.clients {
		select {
		case client.send <- m:
		default:
			// failed to send
			delete(r.clients, client)
			client.Close()
		}
	}
}

func (r *Room) run() {
	log.Printf("run room(id: %d)", r.Id)
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
		case msg := <-r.forward:
			r.broadcast(msg)
		}
	}
}
