package main

import (
	"encoding/binary"
	"fmt"
	"io"
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
	err := write_all(conn, wbuf, int(4+msg_len))

	if err != nil {
		return err
	}

	rbuf := make([]byte, 4+k_max_msg+1)
	err = read_full(conn, rbuf, 4)
	if err != nil {
		return err
	}

	msg_len = binary.LittleEndian.Uint32(rbuf[:4])

	if msg_len > k_max_msg {
		return fmt.Errorf("Message too long!")
	}

	// reply body
	err = read_full(conn, rbuf[4:], int(msg_len))
	if err != nil {
		return err
	}

	// do something
	rbuf[4+msg_len] = '\x00'

	fmt.Println("server says: ", string(rbuf[4:]))

	return nil
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
			panic("should never happen right?")
		}

		n -= rv
		pos += rv
	}
	return nil
}
