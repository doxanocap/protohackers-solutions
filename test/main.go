package main

// import necessary packages
import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	//fmt.Println((65*10 + 71*10 + 20*69 + 89*10 + 88*5 + 25*80 + 10*100 + 5*100) / 100)
	num1 := int64(101)
 
	hex_num := strconv.FormatInt(num1, 16)
	bytes := make([]byte,4)
	bytes[3] = hex_num	
	fmt.Println("hexadecimal num: ", hex_num)
 
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
