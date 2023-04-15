package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func ConnectToServer () net.Conn {
	for {
		var address string
		var port int
		//fmt.Println("Enter IP address:")
		//fmt.Scan(&address)
		//fmt.Scanln()
		address = "localhost"
		for {
			var input string
			fmt.Println("Enter port number:")
			fmt.Scanln(&input)
			num, err := strconv.Atoi(input)
			if err != nil || num < 0 || num > 65535{
				fmt.Println("Error: invalid port number")
				continue
			}
			port = num
			break
		}
		// connection setup
		conn, err := net.DialTimeout("tcp", address+":"+ strconv.Itoa(port), 15 * 1000000000) // 15 sec timeout
		if err != nil {
			fmt.Println(err)
			fmt.Println("Try again")
			continue
		}
		return conn
	}
}

func main() {

	conn := ConnectToServer()
	fmt.Println("You are successfully connected to the server")
	defer conn.Close()

	// goroutine for recieving data from server
	go func() {
		for {
			message := make([]byte, 1024)
			length, err := conn.Read(message)
			if err != nil {
				fmt.Println(err)
				return
			}
			if length > 0 {
				fmt.Print(string(message[:length]))
			}

		}
	}()

	// sending data to server
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
