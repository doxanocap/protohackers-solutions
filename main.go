package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var once sync.Once

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}
	
	tcp, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Couldn't estabilish connection %s", err)
	}

	log.Printf("Listening on port: %s \n", port)

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("Error in accepting connection %s", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	address := conn.RemoteAddr()
	log.Printf("Accepted from %s \n", address)

	socket, err := net.Dial("tcp", "206.189.113.124:16963")
	if err != nil {
		log.Printf("Connecting error %s \n", err)
	}
	once = sync.Once{}
	go ProxyChain(conn, socket)
	ProxyChain(socket, conn)
}

func ProxyChain(first net.Conn, second net.Conn) {		
	defer once.Do(func() { first.Close(); second.Close(); log.Printf("closed connection:") })

	scanner := bufio.NewScanner(first) 
	
	for scanner.Scan() {
		msg  := scanner.Text()

		words := strings.Split(msg, " ")
		newMsg := ""

		for _,v := range words {	
			if len(v) == 0 {
				continue
			}
			if v[0] == '7' && len(v) > 25 && len(v) < 36 && isBoguscoin([]byte(v)) {
				newMsg += "7YWHMfk9JZe0LM0g1ZauHuiSxhI "
				continue
			}
			newMsg += v + " "
		}

		log.Printf("From -{%s}- to -{%s}- | -{Message: %s }- ", first.RemoteAddr(), second.RemoteAddr(), newMsg)
		if _, err := second.Write([]byte(newMsg[:len(newMsg)-1]+"\n")); err != nil {
			log.Printf("Upstream writing error: %s \n",err)
		}
	}
}

func isBoguscoin(data []byte) bool {
	isValid := true
	for _, v := range data {
		if v < 48 || (v > 57 && v < 65) || (v > 90 && v < 97) || v > 122 {
			isValid = false
			break
		}
	}
	return isValid
}

