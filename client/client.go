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
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Подключение к серверу успешно установлено.")

	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Print(message)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Введите сообщение: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = fmt.Fprintf(conn, text)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

