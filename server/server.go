package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

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

type clientInf struct {
	name string
	conn net.Conn
}

func main() {
	var port int
	
	for {
		var input string
        fmt.Println("Enter port number:")
        fmt.Scanln(&input)
        num, err := strconv.Atoi(input)
        if err != nil || num < 0 || num <= 65535{
            fmt.Println("Error: invalid port number")
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

	connectedClients := make(map[clientInf]bool)
	messages := make(chan string)

	// sending messages to clients
	go func() {
		for {
			message := <-messages
			for client := range connectedClients {
				_, err := client.conn.Write([]byte(message))
				if err != nil {
					delete(connectedClients, client)
				}
			}
		}
	}()

	for {
		// accepting new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(conn net.Conn) {
			// creating new user
			client := clientInf{"", conn}
			conn.Write([]byte("Enter your name: "))
			name := make([]byte, 1024)
			for {
				length, err := conn.Read(name)
				if err != nil {
					fmt.Println(err)
					return
				}
				if length > 2 {
					client.name = string(name[:length-2])
					fmt.Println("Connecting the new user:", client.name)
					break
				}
			}
			connectedClients[client] = true
			conn.Write([]byte("You are in chat, enter a message\n"))
			for {
				conn := client.conn
				message := make([]byte, 1024)
				length, err := conn.Read(message)
				if err != nil {
					fmt.Println("User disconnected:", client.name)
					delete(connectedClients, client)
					return
				}
				if length > 0 {			
					messages <- fmt.Sprintf("%s: %s", client.name, message[:length])
					fmt.Printf("%s: %s", client.name, message[:length])
				}
			}
		}(conn)
	}
}
