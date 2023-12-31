package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect: ", err)
		return
	}

	defer conn.Close()

	msg := []byte("hello")
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal("Failed to write: ", err)
		return
	}

	rbuf := make([]byte, 64)
	_, err = conn.Read(rbuf)
	if err != nil {
		log.Fatal("Failed to read: ", err)
		return
	}

	fmt.Println("server says: ", string(rbuf))
}
