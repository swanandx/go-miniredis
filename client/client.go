package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	defer conn.Close()

	err = query(conn, "hello1")
	if err != nil {
		log.Fatal("query failed: ", err)
	}

	err = query(conn, "hello2")
	if err != nil {
		log.Fatal("query failed: ", err)
	}

	err = query(conn, "hello3")
	if err != nil {
		log.Fatal("query failed: ", err)
	}
}

var k_max_msg uint32 = 4096

func query(conn net.Conn, text string) error {
	msg_len := uint32(len(text))
	if msg_len > k_max_msg {
		return fmt.Errorf("Message too long")
	}

	wbuf := make([]byte, 4+k_max_msg)

	binary.LittleEndian.PutUint32(wbuf, msg_len)
	copy(wbuf[4:], text)
	written, err := conn.Write(wbuf[:4+msg_len])

	if err != nil {
		return err
	}

	if written < int(msg_len) {
		return fmt.Errorf("Unable to write full message, wrote %v/%v", written, msg_len)
	}

	rbuf := make([]byte, 4+k_max_msg+1)
	read, err := conn.Read(rbuf)

	if err != nil {
		return err
	}

	if read < 4 {
		return fmt.Errorf("Insufficient bytes, expected at least 4, got %v", read)
	}

	msg_len = binary.LittleEndian.Uint32(rbuf[:4])

	if msg_len > k_max_msg {
		return fmt.Errorf("Message too long!")
	}

	if read < int(msg_len)+4 {
		return fmt.Errorf("Unable to read full message, read %v/%v", read, msg_len)
	}

	// do something
	rbuf[4+msg_len] = '\x00'

	fmt.Println("server says: ", string(rbuf[4:]))

	return nil
}
