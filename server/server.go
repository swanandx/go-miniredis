package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	socket, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal("Failed to open socket: ", err)
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection: ", err)
			continue
		}
		defer conn.Close()
		do_something(conn)
	}
}

func do_something(conn net.Conn) {
	rbuf := make([]byte, 64)
	_, err := conn.Read(rbuf)
	if err != nil {
		log.Fatal("Failed to read: ", err)
		return
	}
	fmt.Println("client says: ", string(rbuf))

	wbuf := []byte("world")
	_, err = conn.Write(wbuf)
	if err != nil {
		log.Fatal("Failed to write: ", err)
		return
	}
}
