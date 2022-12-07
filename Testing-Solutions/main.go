package main

// import necessary packages
import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// var strs = []string{"foo1=bar", "foo2=bar=eqwe", "=ewe","foo3==="}

	// for  _,v := range strs {
	// 	wg.Add(1)
	// 	go func(v string) {
	// 		sendUDPRequest(v)
	// 		wg.Done()
	// 	}(v)
	// }

	var req = []string{"foo1=qwee", "foo2=wqeqe", "version", "foo1"}
	
	for  _,v := range req {
		wg.Add(1)
		go func(v string) {
			sendUDPRequest(v)
			wg.Done()
		}(v)
	}
	wg.Wait()

}


// func sendTCPRequest(i int) {
// 	str := strconv.Itoa(i)
// 	httpRequest := `{"number":` + str + `,"method":"isPrime"}` + "\n"

// 	conn, err := net.Dial("tcp", "localhost:5000")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	if _, err = conn.Write([]byte(httpRequest)); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Print(i, ": ")
// 	if _, err := io.Copy(os.Stdout, conn); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func sendUDPRequest(msg string) {
	conn, err := net.Dial("udp", "169.155.58.11:5000")
	if err != nil {
		log.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(msg)); err != nil {
		fmt.Printf("Write err %v", err)
		os.Exit(-1)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Printf("Read err %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("%s\n", string(buff[:n]))
}