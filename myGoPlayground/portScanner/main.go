package main
import (
 "fmt"
 "net"
)
func main() {
 for i := 1; i <= 10000; i++ {
 	go func(i int) {
 		address := fmt.Sprintf("104.16.244.78:%d", i)
 		conn, err := net.Dial("tcp", address)
 		if err != nil {
 			return
 		}
 		conn.Close()
 		fmt.Printf("%d open\n", i)
 		}(i)
 	}
}
