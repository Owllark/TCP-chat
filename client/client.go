package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	var address string
	var port int
	fmt.Println("Enter IP address:")
	fmt.Scan(&address)
	for {
		var input string
        fmt.Println("Enter port number:")
        fmt.Scanln(&input)
        num, err := strconv.Atoi(input)
        if err != nil || num < 0{
            fmt.Println("Error: invalid port number")
            continue
        }
		port = num
        break
    }
	fmt.Scanln()
	// connecting with server
	conn, err := net.Dial("tcp", address+":"+ strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
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
