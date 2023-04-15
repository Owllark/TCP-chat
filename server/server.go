package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

type serverRequestsType int
type clientRequestsType int

const (
	ENTER_NICKNAME serverRequestsType = iota
	CHOOSE_ROOM
	ENTER_ROOM_PASSWORD
	SEND_MESSAGE
)

const (
	ESCAPE_ROOM serverRequestsType = iota
	GET_USER_LIST
)

type clientState int

const (
	IN_LOBBY clientState = iota
	CREATING_ROOM
	IN_ROOM
)

type clientInf struct {
	name string
	conn net.Conn
	state clientState
}

type roomInf struct {
	name string
	password string
	admin *clientInf
	users map[*clientInf]bool
}

func isPortInUse(port int) bool {
	cmd := exec.Command("cmd", "/C", "netstat -a -n -o")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%d", port)){
			return true
		}
	}

	return false
}

func SendMessage(conn net.Conn, message []byte) (error) {
	_, err := conn.Write([]byte(message))
	return err
}

func GetMessage(conn net.Conn) ([]byte, error) {
	response := make([]byte, 1024)
	for {
		length, err := conn.Read(response)
		if err != nil {
			return []byte(""), err
		}
		if length > 0 {
			return response[:length], err
		}
	}
}

func AcceptClient(conn net.Conn) (clientInf, error) {
	client := clientInf{"", conn, IN_LOBBY}
	SendMessage(client.conn, []byte("Enter your name: "))
	name := make([]byte, 1024)
	for {
		length, err := conn.Read(name)
		if err != nil {
			return client, err
		}
		if length > 2 {
			client.name = string(name[:length-2])
			break
		}
	}
	return client, nil
}

func LaunchRoom(room *roomInf, newClients chan *clientInf, disconnectedClients chan *clientInf) {
	messages := make(chan string, 100)
	// sending messages to clients
	go func() {
		for {
			var message string
			message = <-messages
			fmt.Println("accept message")
			for client := range room.users {
				err := SendMessage(client.conn, []byte(message))
				fmt.Println(message)
				if err != nil {
					delete(room.users, client)
					fmt.Println(err)
				}
			}
		}
	}()
	for {
		var newClient *clientInf
		newClient = <- newClients
		fmt.Println(*newClient)
		room.users[newClient] = true
		go func(client clientInf) {
			for {
				message := make([]byte, 1024)
				message, err := GetMessage(client.conn)
				if err != nil {
					fmt.Println("User disconnected:", client.name)
					delete(room.users, &client)
					disconnectedClients <- &client
					return
				}
				if len(message) > 0 {			
					messages <- fmt.Sprintf("%s: %s", client.name, message)
				}
			}
		} (*newClient)
	}
}

func CreateRoom(client *clientInf) *roomInf{
	var room roomInf
	room.admin = client
	room.users = make(map[*clientInf]bool)
	SendMessage(client.conn, []byte("Enter name of your room:\n"))
	name, _ := GetMessage(client.conn)
	room.name = string(name)
	SendMessage(client.conn, []byte("Enter password of your room or press \"Enter\" to make the room open:\n"))
	password, _ := GetMessage(client.conn)
	room.password = string(password)
	SendMessage(client.conn, []byte("You successfully created the room!\n"))
	return &room
}


func main() {
	var port int
	
	for {
		var input string
        fmt.Println("Enter port number:")
        fmt.Scanln(&input)
        num, err := strconv.Atoi(input)
        if err != nil || num < 0 || num > 65535{
            fmt.Println("Error: invalid port number", err, num)
            continue
        }
		port = num
		if isPortInUse(port) {
			fmt.Printf("Port %d is already in use\n", port)
			continue
		}
        break
    }

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Error listening on the port: %s\n", err)
	}
	defer listener.Close()
	fmt.Println("Server is running and listening on the port " + strconv.Itoa(port))

	connectedClients := make(map[*clientInf]bool)
	rooms := make(map[*roomInf]bool)
	newClients := make(map[*roomInf] chan *clientInf)
	disconnectedClients := make(map[*roomInf] chan *clientInf)


	for {
		// accepting new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(conn net.Conn) {
			// creating new user
			client, err := AcceptClient(conn)
			if err != nil {
				return
			}
			fmt.Println("Connected new user:", client.name)
			connectedClients[&client] = true
			SendMessage(client.conn, []byte("Enter number of room, or enter \"0\" to create a new room:\n"))
			roomList := make([]*roomInf, len(rooms))
			i := 0
			for room := range rooms{
				SendMessage(client.conn, []byte(strconv.Itoa(i + 1) + ". " + room.name))
				roomList[i] = room
				i++
			}
			response, _  := GetMessage(client.conn)
			num, _ := strconv.Atoi(string(response[:len(response) - 2]))
			if num == 0 {
				newRoom := CreateRoom(&client)
				rooms[newRoom] = true
				newClients[newRoom] = make(chan *clientInf, 100)
				disconnectedClients[newRoom] = make(chan *clientInf, 100)
				
				go LaunchRoom(newRoom, newClients[newRoom], disconnectedClients[newRoom])
				newClients[newRoom] <- &client
				fmt.Println("hello, world!")
				
			} else {
				SendMessage(client.conn, []byte("Enter password:\n"))
				response, _  = GetMessage(client.conn)
				if string(response) == roomList[num - 1].password {
					newClients[roomList[num - 1]] <- &client
				} else {
					SendMessage(client.conn, []byte("Wrong password!\n Try again:\n"))
				}
			}
			
		}(conn)
	}
}

