package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

//type tcpConn net.Conn // invalid receiver type *tcpConn (tcpConn is an interface type)
type tcpConn struct {
	net.Conn
}

// {"MessageType":1,"RoomId":1,"Content":"enter"}
func (c *tcpConn) ReadChatMessage(msg *Message) error {
	var bufferSize = 1024
	data := make([]byte, bufferSize)
	n, err := c.Read(data)
	if err != nil {
		log.Println("tcp read error: ", err)
		return err
	}
	log.Println("read data from tcp socket:", data[:n], "EOF")

	json.Unmarshal(data[:n], msg)
	log.Println("read message from tcp socket:", msg)
	msg.CreatedAt = time.Now()
	return nil
}

func (c *tcpConn) WriteChatMessage(msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("message marshal error: ", err)
	}
	_, err = c.Write(b)
	if err != nil {
		log.Println("tcp write error: ", err)
		return err
	}
	return nil
}

func (c *tcpConn) Close() error {
	err := c.Close()
	if err != nil {
		log.Println("tcp close error: ", err)
		return err
	}
	return nil
}

func (c *tcpConn) String() string {
	return c.RemoteAddr().String()
}

func handleTcp(c net.Conn) {
	tcpConn := &tcpConn{Conn: c}
	newClient(tcpConn)
}

func runTcp(addr string) {
	log.Println("start tcp server")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("tcp accept error: ", err)
		}
		handleTcp(c)
	}
}
