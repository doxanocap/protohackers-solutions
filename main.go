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
	Number int `json:"number"`
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
			fmt.Println(err.Error())
			return
		}

		go func(conn net.Conn) {
			defer conn.Close()
			fmt.Println("Connection from:",conn.RemoteAddr())
			
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var raw Prime			
			err = json.Unmarshal([]byte(data),&raw)
			if err != nil || raw.Method != "isPrime" || raw.Number <= 0 {
				fmt.Println(err.Error())
				conn.Write([]byte(data))
				conn.Close()
				return
			}

			isPrime := true
			for i := 2; i < raw.Number; i++ {
				if raw.Number % i == 0 {
					isPrime = false
					break
				}
			}

			res, err := json.Marshal(Result{Method:"isPrime", Prime: isPrime}) 
			
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			conn.Write(res)

		}(conn)
	}
}