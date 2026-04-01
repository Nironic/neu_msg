package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func recv(conn net.Conn) string {
	reader := bufio.NewReader(conn)

	// Читаем до символа \n
	message, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			println("Error EOF")
		}
		return "error"
	}

	// Удаляем \n и пробелы
	text := strings.TrimSpace(message)

	if text == "" {
		return "error"
	}
	return text
}

func send(conn net.Conn, text string) bool {
	_, err := conn.Write([]byte(text + "\n"))
	if err != nil {
		println("Не могу отправить сообщение")
		return false
	}
	return true
}

func server_poll(conn net.Conn, msg string) { // Основная обработка сообщений сервера

}

func polling(conn net.Conn) {
	defer conn.Close()
	// Общий канал отправки сообщений
	tunnel := make(chan string)
	defer close(tunnel)
	go polling_recv(conn, tunnel)

	for {
		select {
		case msg, ok := <-tunnel:
			if !ok {
				return
			}
			println(msg)
		default:
			// SEND GROUP <id_group> <user> <message>
			send(conn, "GET GROUP 1")
			send(conn, "SEND GROUP 1 Roman Привет сосед я сосал табурет")
			send(conn, "GET GROUP 1")
		}
	}

	/*
		for {
			msg, ok := <-tunnel
			if !ok {
				return
			}
			server_poll(conn, msg)
		}
	*/
}

func polling_recv(conn net.Conn, tunnel chan<- string) {
	//Общий канал получения сообщений
	for {
		msg := recv(conn)
		if msg == "error" {
			println("Ошибка при получении сообщения")
			return
		}
		tunnel <- msg
	}
}

func login(conn net.Conn, login string, password string) {
	send(conn, "login")
	if recv(conn) == "OK" {
		send(conn, login)
		recv(conn)
		send(conn, password)
		recv(conn)
		send(conn, "check")
		println(recv(conn))
		polling(conn)
	}
}

func reg(conn net.Conn, login string, password string, username string) {
	send(conn, "reg")
	if recv(conn) == "OK" {
		send(conn, login)
		recv(conn)
		send(conn, password)
		recv(conn)
		send(conn, username)
		println(recv(conn)) // Ответ тут
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	key := "9S2oPsZJ1ipUxKlbyJvr"
	send(conn, key)
	data := recv(conn)
	if data == "KEY" {
		println("Ключ принят")
	} else {
		println("Ключ не принят")
		return
	}
	login(conn, "roma123", "1234")
	//reg(conn, "nk_use", "r02e07m76p76", "Роман Левкин")
}
