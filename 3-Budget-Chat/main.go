package main

// import necessary packages
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	Username string
	Conn     net.Conn
	Error    chan error
}

type Pool struct {
	Register   chan Client
	Unregister chan Client
	Clients    map[Client]bool
	Messages   chan Message
}

type Message struct {
	Sender Client
	Text   string
}

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}

	tcp, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("can't listen on %s/tcp: %s", port, err)
		return
	}

	log.Println("Listening on port: ", port)

	pool := Server()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("conn: %s \n", err)
			return
		}
		go handleConnection(conn, pool)
	}
}


func Start(pool *Pool) {
	for {
		select {
		case client := <-pool.Register:
			if _, ok := pool.Clients[client]; ok {
				client.Error <- fmt.Errorf("user - %s already registered", client.Username)
				break
			}
			if len(client.Username) < 1 || !checkUsername([]byte(client.Username)) {
				client.Error <- fmt.Errorf("invalid username: %s", client.Username)
				break
			}

			content := "*" + client.Username + " has entered the room \n"
			
			log.Println("| --- ENTERED:",content)

			users := []string{}
			for client, online := range pool.Clients {
				if !online {
					continue
				}
				users = append(users, client.Username)
				if _, err := client.Conn.Write([]byte(content)); err != nil {
					log.Println("| --- ERROR -{ registering client }- ", err, client.Conn.RemoteAddr())
				}
			}

			pool.Clients[client] = true
			client.Error <- nil

			content = "* The room contains: " + strings.Join(users, ", ") + "\n"
			if _, err := client.Conn.Write([]byte(content)); err != nil {
				log.Println("| --- ERROR -{ room attendance }- ", err, client.Conn.RemoteAddr())
			}
			log.Println("| --- ALL USERS: ",users, "--- |")

		case client := <-pool.Unregister:
			pool.Clients[client] = false
			content := "*" + client.Username + " has left the room \n"
			
			log.Println("| --- LEFT:",content)
			for client, online := range pool.Clients {
				if !online {
					continue
				}
				if _, err := client.Conn.Write([]byte(content)); err != nil {
					log.Println("| --- ERROR -{ unregistering client }- ", err, client.Conn.RemoteAddr())
				}
			}

		case message := <-pool.Messages:
			if !pool.Clients[message.Sender] {
				break
			}
			sender := message.Sender.Username
			content := "[" + sender + "] " + message.Text + "\n"

			log.Println("|###----- MESSAGE:",content)
			for client, online := range pool.Clients {
				if client.Username == sender || !online {
					continue
				}
				if _, err := client.Conn.Write([]byte(content)); err != nil {
					log.Println("| --- ERROR -{ error while writing messages }- ", err, client.Conn.RemoteAddr())
				}
			}
		}
	}
}

func Server() *Pool {
	newPool := &Pool{
		Register:   make(chan Client),
		Unregister: make(chan Client),
		Clients:    map[Client]bool{},
		Messages:   make(chan Message),
	}
	go Start(newPool)
	return newPool
}


func handleConnection(conn net.Conn, pool *Pool) {
	address := conn.RemoteAddr()

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("| --- ERROR --- Closing error")
		}
	}(conn)

	startMsg := "Welcome to Budget Chat! What shall I call you? \n"
	if _, err := conn.Write([]byte(startMsg)); err != nil {
		log.Println("| --- ERROR --- ", err)
	}

	client := Client{}
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		err := make(chan error)
		client = Client{scanner.Text(), conn, err}
		log.Println("| --- Connected to the remote address:", address, scanner.Text())
		
		pool.Register <- client
		if chanErr := <-err; chanErr != nil {
			log.Println("| --- ERROR -{ chan err }- ", chanErr)
			return
		}

		defer func(conn net.Conn) {
			pool.Unregister <- client
		}(conn)
	} else {
		log.Println("| --- ERROR -{ no data were handled}-")
	}

	for scanner.Scan() {
		pool.Messages <- Message{client, scanner.Text()}
	}
}

func checkUsername(data []byte) bool {
	isValid := true
	for _, v := range data {
		if v < 48 || (v > 57 && v < 65) || (v > 90 && v < 97) || v > 122 {
			isValid = false
			break
		}
	}
	return isValid
}
