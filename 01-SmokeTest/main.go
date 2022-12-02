package main

// import necessary packages
import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	port := "8080"
	if v 	:= os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}

	tcp, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("can't listen on %d/tcp: %s", port, err)
	}

	fmt.Println("listening on port: ", port)

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		go func(conn net.Conn) {
			_, err = io.Copy(conn, conn)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(3)
			}
			fmt.Println(conn.RemoteAddr())
			defer conn.Close()
		}(conn)
	}
}
