package main

// import necessary packages
import (
	"fmt"
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

	tcp, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Printf("can't listen on %s/tcp: %s", port, err)
		return
	}
	fmt.Println("Listening on port: ", port)

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Printf("conn: %s \n", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for i := 0; true; i++ {

	}
}
