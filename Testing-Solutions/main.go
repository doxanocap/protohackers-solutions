package main

// import necessary packages
import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {


	// var wg sync.WaitGroup
	// for i := 0; i < 10000; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		sendRequest(i)
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()
}


func sendRequest(i int) {
	str := strconv.Itoa(i)
	httpRequest := `{"number":` + str + `,"method":"isPrime"}` + "\n"

	conn, err := net.Dial("tcp", "149.248.208.178:5000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(httpRequest)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(i, ": ")
	if _, err := io.Copy(os.Stdout, conn); err != nil {
		fmt.Println(err)
		return
	}
}
