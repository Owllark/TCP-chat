package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	var address string
	var port string
	fmt.Println("Введите ip адрес:")
	fmt.Scan(&address)
	fmt.Println("Введите номер порта:")
	fmt.Scan(&port)
	fmt.Scanln()
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			message := make([]byte, 1024)
			length, err := conn.Read(message)
			if err != nil {
				fmt.Println("ERROR", err)
				return
			}
			if length > 0 {
				fmt.Print(string(message[:length]))
			}

		}
	}()

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
