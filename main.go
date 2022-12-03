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
			fmt.Println("Connection from:",conn.RemoteAddr())
			
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("reading:",err.Error())
				return
			}

			var raw Prime			
			
			if err = json.Unmarshal([]byte(data),&raw); err != nil || raw.Method != "isPrime" || *raw.Number <= 0 || raw.Number == nil {
				fmt.Println(err.Error())
				conn.Write([]byte(data))
				conn.Close()
				return
			}
			isWholeNumber := *raw.Number == float64(int(*raw.Number))

			isPrime := true
			for i := 2; i < int(*raw.Number); i++ {
				if int(*raw.Number) % i == 0 {
					isPrime = false
					break
				}
			}

			res, err := json.Marshal(Result{Method:"isPrime", Prime: isPrime}) 
			
			if err != nil || !isWholeNumber {
				fmt.Println(err.Error())
				conn.Write([]byte("invalid"))
				conn.Close()
				return
			}

			conn.Write(res)
			conn.Close()
		}(conn)
	}
}