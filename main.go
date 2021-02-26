package main

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	conn []net.Conn
}

var connChan = make(chan net.Conn, 40096)

var server Server

func main() {

	lstn, err := net.Listen("tcp", ":2121")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lstn.Accept()
		if err != nil {
			continue
		}

		connChan <- conn

		go addConn(conn)

		go handlerConn(conn)

		go handlerMsgs(conn)

	}
}

func addConn(conn net.Conn) {

	server.conn = append(server.conn, <-connChan)

}

func handlerMsgs(conn net.Conn) {
	for {
		var msg = make([]byte, 40096)
		n, err := conn.Read(msg)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			continue

		}

		if n > 0 {

			for _, c := range server.conn {
				if c == conn {
					continue
				}

				_, err = c.Write(msg[:n])
				if err != nil {
					if err != io.EOF {
						continue
					}
					continue
				}
			}
		}

	}
}

func handlerConn(conn net.Conn) {
	_, err := conn.Write([]byte("You are connected"))
	if err != nil {
		panic(err)
	}
}
