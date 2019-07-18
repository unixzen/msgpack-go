package main

import (
	"github.com/vmihailenco/msgpack"
	"log"
	"net"
	"time"
)

type Message struct {
	Domain string
	Ip     string
}

func main() {

	// Create connections to server
	conn, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		log.Fatal(err)
	}
	conn1, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		log.Fatal(err)
	}

	// Encoding msgpack messages
	message1, err := msgpack.Marshal(&Message{Domain: "yandex.ru", Ip: "192.168.1.10"})
	if err != nil {
		log.Fatal(err)
	}

	message2, err1 := msgpack.Marshal(&Message{Domain: "google.com", Ip: "192.168.1.13"})
	if err1 != nil {
		log.Fatal(err1)
	}

	message3, err := msgpack.Marshal(&Message{Domain: "rambler.ru", Ip: "192.168.1.18"})
	if err != nil {
		log.Fatal(err)
	}

	message4, err1 := msgpack.Marshal(&Message{Domain: "mail.ru", Ip: "192.168.1.44"})
	if err1 != nil {
		log.Fatal(err1)
	}

	// Write to connection
	conn.Write([]byte(message1))
	conn.Write([]byte(message2))
	time.Sleep(4 * time.Second)

	conn1.Write([]byte(message3))
	conn1.Write([]byte(message4))
	conn.Close()
	conn1.Close()

}
