package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Sender string `json:"sender"`
	Body   string `json:"body"`
}

var msgChan = make(chan Message, 40096)

var messageMe Message
var messageFromServer Message

func main() {
	var sender string
	fmt.Printf("Your name: ")
	fmt.Scanf("%s\n", &sender)

	messageMe.Sender = sender

	conn, err := net.Dial("tcp", ":2121")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	for {

		go readMsg(conn)

		go logMsg()

		sendMsg(conn)

	}
}

func logMsg() {
	var msg Message

	for {
		msg = <-msgChan

		if len(msg.Body) > 0 {
			fmt.Printf("\n[ %s ] - said: %s\n", msg.Sender, msg.Body)
		}
	}
}

func sendMsg(conn net.Conn) {
	var msg string
	var sendMsg = make([]byte, 40096)
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err = reader.ReadString('\n')
		if err != nil {
			continue
		}

		messageMe.Body = msg

		sendMsg, err = json.Marshal(messageMe)

		_, err = conn.Write(sendMsg)
		if err != nil {
			continue
		}
	}
}

func readMsg(conn net.Conn) {
	var msg = make([]byte, 40096)
	for {

		n, err := conn.Read(msg)
		if err != nil {
			continue
		}

		if n > 0 {
			err = json.Unmarshal(msg[:n], &messageFromServer)
			if err != nil {
				continue
			}

			msgChan <- messageFromServer

		}
	}
}
