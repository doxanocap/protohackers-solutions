package main

// import necessary packages
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Result struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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
	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error:", err)
			if err == io.EOF {
				break
			}
		}

		fmt.Printf("Incomming request: %s \n", req)

		var raw = &Request{}

		if err = json.Unmarshal([]byte(req), &raw); err != nil || raw.Method != "isPrime" || raw.Number == nil {
			fmt.Printf("| Unmarshalling err -> Malformend: %s \n", err)
			if _, err := conn.Write([]byte("invalid" + "\n")); err != nil {
				fmt.Printf("Writing err: %s", err)
			}
			break
		}

		isWholeNumber := *raw.Number == float64(int(*raw.Number))

		prime := isWholeNumber && isPrime(int(*raw.Number))

		res, err := json.Marshal(Result{Method: "isPrime", Prime: prime})
		if err != nil {
			fmt.Printf("| Marshaling error: %s\n", err)
		}

		res = append(res, '\n')
		if _, err := conn.Write(res); err != nil {
			fmt.Printf("| Connection writing error: %s", err)
		}

	}

	// FAIL:830818 is composite but response said prime

}

func isPrime(n int) bool {
	if n == 2 || n == 3 {
		return true
	}

	if n <= 1 || n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}
