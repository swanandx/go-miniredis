package main

import (
	"encoding/binary"
	"fmt"
	"io"
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
			log.Println("Failed to accept connection: ", err)
			continue
		}
		defer conn.Close()

		for {
			e := one_request(conn)
			if e != nil {
				log.Println("error while processing requests: ", e)
				break
			}
		}
	}
}

var k_max_msg uint32 = 4096

func one_request(conn net.Conn) error {
	// 4 bytes header
	rbuf := make([]byte, 4+k_max_msg+1)
	err := read_full(conn, rbuf, 4)
	if err != nil {
		return err
	}

	msg_len := binary.LittleEndian.Uint32(rbuf[:4])

	if msg_len > k_max_msg {
		return fmt.Errorf("Message too long!")
	}

	// request body
	err = read_full(conn, rbuf[4:], int(msg_len))
	if err != nil {
		return err
	}

	// do something
	rbuf[4+msg_len] = '\x00'

	fmt.Println("client says: ", string(rbuf[4:]))

	// reply using same protocol
	reply := []byte("world")
	msg_len = uint32(len(reply))
	wbuf := make([]byte, 4+msg_len)
	binary.LittleEndian.PutUint32(wbuf, msg_len)
	copy(wbuf[4:], reply)
	return write_all(conn, wbuf, int(4+msg_len))
}

func read_full(c net.Conn, buf []byte, n int) error {
	reader := io.LimitReader(c, int64(n))
	pos := 0
	for n > 0 {
		rv, err := reader.Read(buf[pos:])
		if err != nil {
			return err
		}

		if rv > n {
			panic("we have limited the reader, we should never reach here")
		}

		n -= rv
		pos += rv
	}

	return nil
}

func write_all(c net.Conn, buf []byte, n int) error {
	pos := 0
	for n > 0 {
		rv, err := c.Write(buf[pos:n])
		if err != nil {
			return err
		}

		if rv > n {
			panic("we should never reach here")
		}

		n -= rv
		pos += rv
	}
	return nil
}
