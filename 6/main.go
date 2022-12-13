package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Plate struct {
	Plate     string
	Mile      uint16
	Timestamp uint32
	Limit     uint16
}

type Camera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

type Ticket struct {
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      float32
}

var Plates map[uint16][]Plate
var Cameras map[uint16][]Camera
var Tickets map[uint16][]Ticket
var Cars map[string]bool
var Counter int

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

	Plates = map[uint16][]Plate{}
	Cameras = map[uint16][]Camera{}
	Tickets = map[uint16][]Ticket{}
	Cars = map[string]bool{}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Printf("Error in accepting connection %s", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	fmt.Println("Client from", conn.RemoteAddr())
	var clientType string
	var camera Camera
	var roads []uint16
	tempCounter := 0

	for {
		buff := make([]byte, 1024)
		if clientType == "Dispatcher" {
			for {
				if Counter > tempCounter {
					fmt.Println(clientType)
					SendMsgDispathcer(roads, conn)
					tempCounter = Counter
				}
			}
		}

		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err)
			break
		}

		buff = buff[:n]
		fmt.Println(buff)
		if len(buff) == 0 {
			continue
		}

		switch buff[0] {
		case 64:
			fmt.Println("64: ", buff[1:5])
			go Heartbeat(conn, buff[1:5])
		case 32:
			fmt.Println("| --- Buffer --- 1", buff)
			for i := 0; i < len(buff); i++ {
				fmt.Println("|---", camera, i)
				if buff[i] == 32 && i < len(buff)-1 {
					ln := int(buff[i+1])
					HandlePlate(camera, buff[i:i+ln+6])
					i += ln + 5
				}
			}

		case 128:
			clientType = "Camera"

			roadId := binary.BigEndian.Uint16(buff[1:3])
			camera = Camera{
				Road:  roadId,
				Mile:  binary.BigEndian.Uint16(buff[3:5]),
				Limit: binary.BigEndian.Uint16(buff[5:7]),
			}

			Cameras[roadId] = append(Cameras[roadId], camera)

			if len(buff) <= 7 {
				continue
			}

			buff = buff[7:]

			fmt.Println("| --- Buffer --- 2", buff)
			for i := 0; i < len(buff); i++ {
				fmt.Println("|---", camera, i)
				if buff[i] == 32 && i < len(buff)-1 {
					ln := int(buff[i+1])
					HandlePlate(camera, buff[i:i+ln+6])
					i += ln + 5
				}
			}

		case 129:
			clientType = "Dispatcher"
			ln := int(buff[1])
			buff = buff[2:]
			for i := 0; i < ln; i++ {
				roads = append(roads, binary.BigEndian.Uint16(buff[2*i:2*i+2]))
			}
			continue
		}
	}
}

func GetTicket(roadId, mile1, mile2 uint16, timestamp1, timestamp2 uint32, plate string, speed float32) Ticket {
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

func HandlePlate(camera Camera, buff []byte) {
	roadId := camera.Road
	ln := int(buff[1])

	carName := string(buff[2 : ln+2])
	timeStamp := binary.BigEndian.Uint32(buff[ln+2 : ln+6])

	plate := Plate{Plate: carName, Timestamp: timeStamp, Mile: camera.Mile, Limit: camera.Limit}
	fmt.Printf("| ---- NEW PLATE ---- ")
	fmt.Print(plate)
	fmt.Println("|")

	Plates[roadId] = append(Plates[roadId], plate)

	maxSpeed := float32(-1.0)
	idx := -1

	for i := 0; i < len(Plates[roadId]); i++ {
		if Plates[roadId][i].Plate != plate.Plate || plate.Timestamp-Plates[roadId][i].Timestamp == 0 {
			continue
		}
		if plate.Timestamp < Plates[roadId][i].Timestamp {
			Plates[roadId][i].Timestamp, plate.Timestamp = plate.Timestamp, Plates[roadId][i].Timestamp
		}
		if plate.Mile < Plates[roadId][i].Mile {
			Plates[roadId][i].Mile, plate.Mile = plate.Mile, Plates[roadId][i].Mile
		}
		speed := (float32(plate.Mile-Plates[roadId][i].Mile) * 3600) / float32(plate.Timestamp-Plates[roadId][i].Timestamp)
		if speed > maxSpeed {
			maxSpeed = speed
			idx = i
		}
	}

	if idx == -1 {
		return
	}

	ticket := GetTicket(roadId, Plates[roadId][idx].Mile, plate.Mile, Plates[roadId][idx].Timestamp, plate.Timestamp, carName, maxSpeed)

	if maxSpeed >= 0.5+float32(camera.Limit) {

		fmt.Println("|")
		fmt.Println("| ----Ticket --- ", ticket)
		fmt.Println("|")
		fmt.Printf("| ------- Car: %s ---- with speed %f while limit is %d \n\n", carName, maxSpeed, camera.Limit)
		fmt.Println("|")
		Tickets[roadId] = append(Tickets[roadId], ticket)
		Counter++
	}
}

func SendMsgDispathcer(roads []uint16, conn net.Conn) {
	tempTicket := map[uint16][]Ticket{}
	for _, roadId := range roads {
		for _, ticket := range Tickets[roadId] {
			if _, ok := Cars[ticket.Plate]; ok {
				msg := []byte{16}
				msg = append(msg, []byte("Error")...)
				if _, err := conn.Write(msg); err != nil {
					fmt.Println("Writing msg1", err)
				}
				tempTicket[roadId] = append(tempTicket[roadId], ticket)
				continue
			} else {
				fmt.Println(ticket.Plate)
			}
			ticket.Speed *= 100
			ticket.Speed = float32(int(ticket.Speed))
			msg := GetTicketBytes(ticket)
			if _, err := conn.Write(msg); err != nil {
				fmt.Println("Writing msg2", err)
			}
			Cars[ticket.Plate] = true
		}
	}
	Tickets = tempTicket
}

func GetTicketBytes(ticket Ticket) []byte {
	msg := []byte{33, byte(len(ticket.Plate))}
	msg = append(msg, []byte(ticket.Plate)...)

	temp, temp2 := make([]byte, 2), make([]byte, 4)
	binary.BigEndian.PutUint16(temp[:], ticket.Road)
	msg = append(msg, temp...)
	temp = make([]byte, 2)
	binary.BigEndian.PutUint16(temp[:], ticket.Mile1)
	msg = append(msg, temp...)
	temp = make([]byte, 2)
	binary.BigEndian.PutUint32(temp2[:], ticket.Timestamp1)
	msg = append(msg, temp2...)
	temp2 = make([]byte, 4)
	binary.BigEndian.PutUint16(temp[:], ticket.Mile2)
	msg = append(msg, temp...)
	temp = make([]byte, 2)
	binary.BigEndian.PutUint32(temp2[:], ticket.Timestamp2)
	msg = append(msg, temp2...)
	temp2 = make([]byte, 4)
	binary.BigEndian.PutUint16(temp[:], uint16(ticket.Speed))
	msg = append(msg, temp...)
	temp = make([]byte, 2)
	return msg
}
