package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
)

func closeit(fd int) {
	// shall we use this instead?
	// syscall.CloseOnExec(fd)
	err := syscall.Close(fd)
	if err != nil {
		log.Fatal("Failed to close the fd: ", err)
	}
}

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Println("Failed to open socket: ", err)
		return
	}

	defer closeit(fd)

	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	addr := syscall.SockaddrInet4{
		Port: 1234,
	}

	err = syscall.Bind(fd, &addr)
	if err != nil {
		log.Println("Failed to bind socket: ", err)
		return
	}

	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		log.Println("Failed to bind socket: ", err)
		return
	}

	for {
		conn, _, err := syscall.Accept(fd)
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}

		for {
			e := one_request(conn)
			if e != nil {
				log.Println("error while processing requests: ", e)
				break
			}
		}

		closeit(conn)
	}
}

var k_max_msg uint32 = 4096

func one_request(conn int) error {
	// 4 bytes header
	rbuf := make([]byte, 4+k_max_msg+1)
	read, err := syscall.Read(conn, rbuf)

	if err != nil {
		return err
	}

	if read == 0 {
		return fmt.Errorf("EOF")
	}

	if read < 4 {
		return fmt.Errorf("Insufficient bytes, expected at least 4, got %v", read)
	}

	msg_len := binary.LittleEndian.Uint32(rbuf[:4])

	if msg_len > k_max_msg {
		return fmt.Errorf("Message too long!")
	}

	if read < int(msg_len)+4 {
		return fmt.Errorf("Unable to read full message, read %v/%v", read, msg_len)
	}

	rbuf[4+msg_len] = '\x00'

	fmt.Println("client says: ", string(rbuf[4:]))

	// reply using same protocol
	reply := []byte("world")
	msg_len = uint32(len(reply))
	wbuf := make([]byte, 4+msg_len)
	binary.LittleEndian.PutUint32(wbuf, msg_len)
	copy(wbuf[4:], reply)

	written, err := syscall.Write(conn, wbuf[:4+msg_len])

	if err != nil {
		return err
	}

	if written < int(msg_len) {
		return fmt.Errorf("Unable to write full message, wrote %v/%v", written, msg_len)
	}

	return nil
}
