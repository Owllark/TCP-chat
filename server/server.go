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
		if strings.Contains(line, fmt.Sprintf(":%d", port)) && strings.Contains(line, "LISTENING") {
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
	fmt.Println("Введите номер порта:")
	fmt.Scan(&port)
	for isPortInUse(port) {
		fmt.Printf("Порт %d уже используется\n", port)
		fmt.Println("Введите номер порта:")
		fmt.Scan(&port)
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Ошибка при прослушивании порта: %s\n", err)
	}
	defer listener.Close()
	fmt.Println("Сервер запущен и слушает порт " + strconv.Itoa(port))

	connectedClients := make(map[clientInf]bool)
	messages := make(chan string)

	go func() {
		for {
			message := <-messages
			for client := range connectedClients {
				_, err := client.conn.Write([]byte(message + "\n"))
				if err != nil {
					delete(connectedClients, client)
				}
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func(conn net.Conn) {
			client := clientInf{"", conn}
			conn.Write([]byte("Введите свое имя: "))
			for {
				name := make([]byte, 1024)
				length, err := conn.Read(name)
				if err != nil {
					return
				}
				if length > 0 {
					client.name = string(name[:length])
					fmt.Println("Подключение нового клиента:", client.name)
					break
				}
			}
			connectedClients[client] = true
			for {
				conn := client.conn
				message := make([]byte, 1024)
				length, err := conn.Read(message)
				if err != nil {
					fmt.Println("Клиент отключился:", client.name)
					delete(connectedClients, client)
					return
				}
				if length > 0 {
					messages <- fmt.Sprintf("%s: %s", client.name, message[:length])
					fmt.Printf("%s: %s", client.name, message[:length])
					fmt.Println(client.name)
				}
			}
		}(conn)
	}
}
