package main

// import necessary packages
import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Insert struct {
	Timestamp int64
	Price     int64
}

type Query struct {
	Mintime int64
	Maxtime int64
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
	
	var char string

	hex_data := []byte{}
	var data1, data2 int64
	history := []Insert{}

	for i := 0; true; i++ {

		buff := make([]byte, 1)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		if i%9 == 0 {
			char = string(buff[:n])
			continue
		}

		hex_data = append(hex_data, buff[:n]...)

		if i%9 == 4 {
			data1, _ = strconv.ParseInt(byteToString(hex_data), 16, 64)
			if char == "Q" {
				fmt.Println("Hexdata", hex_data)
			}
			hex_data = nil
		}

		if i%9 == 8 {
			data2, _ = strconv.ParseInt(byteToString(hex_data), 16, 64)
			if char == "Q" {
				fmt.Println("Hexdata", hex_data)
			}
			if char == "I" {
				history = append(history, Insert{data1, data2})
			} else if char == "Q" {
				n, total := int64(0), int64(0)
				fmt.Println("| - Range :",data1, data2, "- | ", history)
				for _, ins := range history {
					if ins.Timestamp >= data1 && ins.Timestamp <= data2 {
						total += ins.Price
						n++
					}
				}

				result_data := total
				if n != 0 {
					result_data = total / n
				}
				
				result_bytes := numToHexToBytes(result_data)

				fmt.Println("| ---- Results: ",result_data, result_bytes, n)
				
				if _, err := conn.Write(result_bytes); err != nil {
					fmt.Printf("Writing error %s \n", err)
				}
				break
			}
			//fmt.Println(char, "| Decoded data:", history)			
			hex_data = nil
		}
	}
}

func byteToString(bytes []byte) string {
	str := ""
	for _, v := range bytes {
		if v < 10 {
			str += "0" + string(v+48)
			continue
		}
		temp := ""
		for v > 0 {
			digit := v % 10
			temp = string(digit+48) + temp
			v = v / 10
		}
		str += temp
	}
	return str
}

func numToHexToBytes(num int64) []byte {
	hex_str := strconv.FormatInt(num, 16)
	hex_num, _ := strconv.Atoi(hex_str)
	bytes := make([]byte,4)

	bytes[3] = byte(hex_num)	
	return bytes
}
