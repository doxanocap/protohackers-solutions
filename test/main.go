package main

// import necessary packages
import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Prime struct {
	Method string `json:"method"`
	Number *float64 `json:"number"`
}

type Result struct {
	Method string `json:"method"`
	Prime bool `json:"prime"`
}

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("can't listen on %s/tcp: %s", port, err)
	}

	fmt.Println("listening on port: ", port)

	for {
		
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Println("tcp err",err.Error())
			return
		}

		go func(conn net.Conn) {
			co
		}(conn)
	}
}