package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Plate struct {
	Plate     string
	Timestamp int
	Ticket    Ticket
}

type Camera struct {
	Road  int
	Mile  int
	Limit int
}

type Ticket struct {
	Plate      string
	Road       int
	Mile1      int
	Timestamp1 int
	Mile2      int
	Timestamp2 int
	Speed      int
}

var Plates map[int][]Plate
var Cameras map[int][]Camera

func main() {
	port := "8080"
	if v := os.Getenv("ECHO_PORT"); v != "" {
		port = v
	}

	tcp, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Couldn't estabilish connection %s", err)
	}

	log.Printf("Listening on port: %s \n", port)

	Plates = map[int][]Plate{}
	Cameras = map[int][]Camera{}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("Error in accepting connection %s", err)
		}

		ch := make(chan string)
		go handleConnection(conn, ch)
	}
}

func handleConnection(conn net.Conn, ch chan string) {
	defer func() {
		conn.Close()
	}()

	clientType := ""
	camera := Camera{}
	roads := []int{}
	carName := ""
	first := true
	fmt.Println("QWEQEW")
	for {
		if clientType == "Dispatcher" {
			if first == true {
				first = false
				for _, roadId := range roads {
					for i, plate := range Plates[roadId] {
						if plate.Ticket.Speed >= Cameras[roadId][i].Limit {
							plate.Ticket.Speed *= 100
							fmt.Println(plate.Ticket)
							msg := []byte{33, byte(len(plate.Ticket.Plate))}
							msg = append(msg, []byte(plate.Ticket.Plate)...)
							i = 2 + len(plate.Ticket.Plate)
							temp := make([]byte, 2)
							binary.BigEndian.PutUint16(temp[:], uint16(plate.Ticket.Road))
							msg = append(msg, temp...)
							temp = make([]byte, 2)

							binary.BigEndian.PutUint16(temp[:], uint16(plate.Ticket.Mile1))
							msg = append(msg, temp...)
							temp = make([]byte, 4)

							binary.BigEndian.PutUint32(temp[:], uint32(plate.Ticket.Timestamp1))
							msg = append(msg, temp...)
							temp = make([]byte, 2)

							binary.BigEndian.PutUint16(temp[:], uint16(plate.Ticket.Mile2))
							msg = append(msg, temp...)
							temp = make([]byte, 4)

							binary.BigEndian.PutUint32(temp[:], uint32(plate.Ticket.Timestamp2))
							msg = append(msg, temp...)
							temp = make([]byte, 2)

							binary.BigEndian.PutUint16(temp[:], uint16(plate.Ticket.Speed))
							msg = append(msg, temp...)
							temp = make([]byte, 2)

							fmt.Println(msg)
							if _, err := conn.Write([]byte(msg)); err != nil {
								fmt.Println("QWEQEWQE", err)
							}
						}
					}
				}
				continue
			}
			select {
			case msg1 := <-ch:
				fmt.Println(msg1)
				if msg1 == "Check" {
					for _, roadId := range roads {
						fmt.Println("rOADID", roadId)
						for i, plate := range Plates[roadId] {
							//fmt.Println(plate.Ticket, Cameras[roadId][i].Limit)
							if plate.Ticket.Speed >= Cameras[roadId][i].Limit {
								plate.Ticket.Speed *= 100
								ticketString, _ := json.Marshal(plate.Ticket)
								msg := "Ticket" + string(ticketString)
								fmt.Println(msg)
								if _, err := conn.Write([]byte(msg)); err != nil {
									fmt.Println("QWEQEWQE", err)
								}
							}
						}
					}
				}
			}
			continue
		}

		buff := make([]byte, 1024)
		time.Sleep(time.Millisecond * 200)
		n, err := conn.Read(buff)

		if err != nil {
			fmt.Println(err)
			break
		}
		buff = buff[:n]

		if buff[0] == 64 {
			fmt.Println("64: ", buff[1:5])
			go Heartbeat(conn, buff[1:5])
		}

		fmt.Println("Buffer:", buff)

		if buff[0] == 128 {
			clientType = "Camera"

			roadId := int(binary.BigEndian.Uint16(buff[1:3]))
			camera = Camera{
				Road:  roadId,
				Mile:  int(binary.BigEndian.Uint16(buff[3:5])),
				Limit: int(binary.BigEndian.Uint16(buff[5:7])),
			}

			prevCam := Camera{}
			prevPlate := Plate{}
			if len(Cameras[roadId]) > 0 {
				prevCam = Cameras[roadId][len(Cameras[roadId])-1]
				prevPlate = Plates[roadId][len(Plates[roadId])-1]
			}

			Cameras[roadId] = append(Cameras[roadId], camera)

			buff = buff[7:]
			ln := int(buff[1])

			carName = string(buff[2 : 2+ln])
			timeStamp := int(binary.BigEndian.Uint32(buff[2+ln:]))

			plate := Plate{Plate: carName, Timestamp: timeStamp, Ticket: Ticket{}}
			if prevCam.Road == 0 || prevPlate.Plate == "" {
				Plates[roadId] = append(Plates[roadId], plate)
				continue
			}

			speed := (camera.Mile - prevCam.Mile) * 3600 / (timeStamp - prevPlate.Timestamp)
			fmt.Println(speed)
			plate.Ticket = GetTicket(roadId, prevCam.Mile, camera.Mile, prevPlate.Timestamp, timeStamp, carName, speed)
			Plates[roadId] = append(Plates[roadId], plate)
			ch <- "Check"
		} else if buff[0] == 129 {
			clientType = "Dispatcher"
			ln := int(buff[1])
			buff = buff[2:]
			for i := 0; i < ln; i++ {
				roads = append(roads, int(binary.BigEndian.Uint16(buff[2*i:2*i+2])))
			}
			continue
		}
	}
}

func GetTicket(roadId, mile1, mile2, timestamp1, timestamp2 int, plate string, speed int) Ticket {
	ticket := Ticket{
		Plate:      plate,
		Road:       roadId,
		Mile1:      mile1,
		Timestamp1: timestamp1,
		Mile2:      mile2,
		Timestamp2: timestamp2,
		Speed:      speed,
	}
	return ticket
}

func Heartbeat(conn net.Conn, data []byte) {
	interval := binary.BigEndian.Uint32(data)
	if interval == 0 {
		return
	}

	for {
		if _, err := conn.Write([]byte{65}); err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Heartbeat")
		time.Sleep(time.Second * time.Duration(interval/10))
	}
}
