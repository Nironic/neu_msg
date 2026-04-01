package main

import (
	"bufio"
	"io"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	myApp := app.New()
	myWindow := myApp.NewWindow("Простой Чат")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Область чата (только для чтения)
	chatDisplay := widget.NewMultiLineEntry()
	chatDisplay.Disable()
	chatDisplay.SetPlaceHolder("История сообщений...")

	// Поле ввода
	messageInput := widget.NewEntry()
	messageInput.SetPlaceHolder("Введите сообщение...")

	// Обработка отправки
	sendMessage := func() {
		message := strings.TrimSpace(messageInput.Text)
		if message != "" {
			// Добавляем сообщение пользователя
			chatDisplay.SetText(chatDisplay.Text + "Вы: " + message + "\n")
			messageInput.SetText("")

			// Ответ бота
			msg := strings.ToLower(message)
			var response string
			if msg == "привет" {
				response = "Привет! Как дела?"
			} else if msg == "пока" {
				response = "До свидания! Было приятно пообщаться."
			} else if strings.Contains(msg, "?") {
				response = "Это интересный вопрос! Я пока не умею на него отвечать."
			} else {
				response = "Вы написали: " + message
			}
			chatDisplay.SetText(chatDisplay.Text + "Бот: " + response + "\n")

			// Прокрутка вниз (автоматически)
		}
	}

	// Отправка по Enter
	messageInput.OnSubmitted = func(string) {
		sendMessage()
	}

	// Кнопка отправки
	sendButton := widget.NewButton("Отправить", sendMessage)

	// Компоновка
	content := container.NewBorder(
		nil, // top
		container.NewBorder(nil, nil, nil, sendButton, messageInput), // bottom
		nil,         // left
		nil,         // right
		chatDisplay, // center
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
