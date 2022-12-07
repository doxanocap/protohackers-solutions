package main

// import necessary packages
import (
	"encoding/binary"
	"fmt"
	"io"
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

	char := []byte{}
	history := []Insert{}
	hex_data := []byte{}
	var data1, data2 int32

	for i := 0; true; i++ {
		buff := make([]byte, 1)
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error:", err)
			break
		}

		if i%9 == 0 {
			char = buff[:n]
			continue
		}

		hex_data = append(hex_data, buff[:n]...)
		if i%9 == 4 {
			data1 = int32(binary.BigEndian.Uint32(hex_data))
			hex_data = nil
		}

		if i%9 == 8 {
			data2 = int32(binary.BigEndian.Uint32(hex_data))
			hex_data = nil

			if char[0] == 'I' {
				history = append(history, Insert{data1, data2})
				continue
			}

			fmt.Println("|-- History :", len(history))
			fmt.Println("|---- Range :", data1, data2, "- | ")

			var n, total, result_data int
			for _, ins := range history {
				if data1 <= ins.Timestamp && ins.Timestamp <= data2 {
					total += int(ins.Price)
					n++
				}
			}

			if n != 0 {
				result_data = total / n
			}

			result_bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(result_bytes, uint32(result_data))

			fmt.Println("| ------ Results: ", result_data, result_bytes, n)
			if _, err := conn.Write(result_bytes); err != nil {
				fmt.Printf("Writing error %s \n", err)
			}
		}
	}
}
