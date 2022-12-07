package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}
	
	conn, err := net.ListenPacket("udp", fmt.Sprintf("fly-global-services:%s", port))
	if err != nil {
		log.Printf("can't listen on %s/udp: %s", port, err)
	}
	
	log.Println(port)
	
	database := map[string]string{
		"version": "Ken's Key-Value Store 1.0",
	}

	buf := make([]byte, 1024)		
	for {
		n,addr,err := conn.ReadFrom(buf)
        if err != nil {
            fmt.Println("Error: ",err)
			break
		}

		request := string(buf[:n])
		log.Printf("\n|--- From %s \n|--- Received %s \n",addr,string(buf[0:n]))


		key, value, ok := strings.Cut(request, "=")
		if ok {
			if key == "version" {
				continue
			}
			database[key] = value

		} else {
			response := fmt.Sprintf("%v=%v", key, database[key])

			if _, err = conn.WriteTo([]byte(response), addr); err != nil {
				log.Println("|--- Writing ----", err, response)
			}
		}
	}	
}
