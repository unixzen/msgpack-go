package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/vmihailenco/msgpack"
	"log"
	"net"
	"time"
)

// Convert Ip string to uint32
func ip2int(ip string) uint32 {
	ipToUint32 := net.ParseIP(ip)
	if len(ipToUint32) == 16 {
		return binary.BigEndian.Uint32(ipToUint32[12:16])
	}
	return binary.BigEndian.Uint32(ipToUint32)
}

// Define server and accept connection
func server(address string, port string, c *cache.Cache) {

	socket := address + ":" + port
	fmt.Printf("Starting server...%s\n", socket)

	ln, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatal(err)

	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go handleRequest(conn, c)
	}
}

// Handle connections at server
func handleRequest(conn net.Conn, c *cache.Cache) {

	// Decode msgpack messages with buffering
	dec := msgpack.NewDecoder(conn)

	for {
		m := map[string]interface{}{}

		err := dec.Decode(&m)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get expiration for messsages
		_, found := c.Get(m["Domain"].(string))

		// If not set expiration of message, set expiration for message
		if !found {
			c.Set(m["Domain"].(string), m["Ip"], cache.DefaultExpiration)
		}
	}
}

// Print messages
func printOut(c *cache.Cache) {
	for {
		if len(c.Items()) == 0 {
			continue
		}
		log.Println("==============================")
		for k, v := range c.Items() {
			ipUint := ip2int(v.Object.(string))
			log.Printf("%s, %d", k, ipUint)
		}
		log.Println("==============================")
		time.Sleep(time.Second)
	}
}

func main() {
	tcpaddr := flag.String("tcpaddr", "", "Listening network address")
	tcpport := flag.String("tcpport", "", "Listening port")
	flag.Parse()

	if *tcpaddr == "" || *tcpport == "" {
		*tcpaddr = "127.0.0.1"
		*tcpport = "6000"
	}

	// Create cache with expiration time for message
	c := cache.New(10*time.Second, 10*time.Second)

	close := make(chan bool)
	go server(*tcpaddr, *tcpport, c)
	go printOut(c)
	<-close

}
