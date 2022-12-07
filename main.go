package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

type Insert struct {
	Timestamp int32
	Price     int32
}

type Query struct {
	Mintime int32
	Maxtime int32
}

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}
	
	tcp, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Couldn't estabilish connection %s", err)
	}
	fmt.Println("Listening on port: ", port)

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("Error in accepting connection %s", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

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

	}
	return isValid
}
