package main

// import necessary packages
import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	//fmt.Println((65*10 + 71*10 + 20*69 + 89*10 + 88*5 + 25*80 + 10*100 + 5*100) / 100)

	// use the FormatInt() function to convert decimal to hexadecimal
	// store the result in output variable
	buf := new(bytes.Buffer)
	var num int64 = 257
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf)

	// str := ""
	// bytes := []byte{0, 0, 39, 49}
	// for _, v := range bytes {
	// 	if v < 10 {
	// 		str += "0" + string(v+48)
	// 		continue
	// 	}
	// 	strByte := ""
	// 	for v > 0 {
	// 		digit := v % 10
	// 		strByte = string(digit+48) + strByte
	// 		v = v / 10
	// 	}
	// 	str += strByte
	// }
	// fmt.Println(str)

	// hex_num = "000003e8"

	// num, err = strconv.ParseInt(hex_num, 16, 64)

	// if err != nil {

	// 	panic(err)

	// }

	// fmt.Println("decimal num: ", num)

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

func hexToDecimal(data []byte) string {
	dst := make([]byte, hex.DecodedLen(len(data)))
	n, err := hex.Decode(dst, data)
	if err != nil {
		fmt.Println("err")
	}

	return string(dst[:n])
}
