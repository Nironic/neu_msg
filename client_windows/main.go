package main

import (
	"bufio"
	"fmt"
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
			println("Ошибка чтения")
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
		fmt.Println("Не могу принять сообщение от клиента")
		return false
	}
	return true
}

func key_login(conn net.Conn, key string) chan string {
	send(conn, key)
	fmt.Println(recv(conn))
	send(conn, "login")
	fmt.Println(recv(conn))
	login := "roma123"
	password := "1234"
	send(conn, login)
	recv(conn)
	send(conn, password)
	recv(conn)
	send(conn, "check")
	fmt.Println(recv(conn))
	tunnel := make(chan string)
	go polling_recv(conn, tunnel)
	return tunnel
}

func polling_recv(conn net.Conn, tunnel chan<- string) {
	//Общий канал получения сообщений
	for {
		msg := recv(conn)
		if msg == "error" {
			fmt.Printf("Ошибка при получении сообщения [%s]", conn.RemoteAddr().String())
			return
		}
		tunnel <- msg
	}
}

func main() {
	key := "9S2oPsZJ1ipUxKlbyJvr"
	a := app.New()
	w := a.NewWindow("Messenger Client")
	w.Resize(fyne.NewSize(400, 600))

	// Подключение к серверу
	conn, err := net.Dial("tcp", "127.0.0.1:8080") // поменяй адрес
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}

	// Чат (только чтение)
	chat := widget.NewMultiLineEntry()
	chat.SetPlaceHolder("Чат...")
	chat.Disable()

	scroll := container.NewVScroll(chat)

	// Поле ввода
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите сообщение...")

	tunnel := key_login(conn, key)
	send(conn, "GET GROUP 1")

	go func() {
		result := ""
		for {
			select {
			case msg := <-tunnel:
				fmt.Println(msg)
				result += msg + "\n"
				//chat.Text = ""
				if msg == "END" {
					fyne.Do(func() {
						chat.SetText(result)
						scroll.ScrollToBottom()
					})
				}
			default:
				//time.Sleep(500 * time.Millisecond)
				//send(conn, "GET GROUP 1")
			}
		}
	}()

	sendMessage := func() {
		text := input.Text
		if text == "" {
			return
		}
		msg := "SEND GROUP 1 roma123 " + text
		ok := send(conn, msg)
		if !ok {
			fmt.Println("Не удалось отправить")
			return
		}

		input.SetText("")
	}
	// Кнопка
	sendBtn := widget.NewButton("Отправить", sendMessage)

	// Enter = отправка
	input.OnSubmitted = func(string) {
		sendMessage()
	}

	// Нижняя панель
	bottom := container.NewBorder(nil, nil, nil, sendBtn, input)

	// Layout
	content := container.NewBorder(nil, bottom, nil, nil, scroll)
	w.SetContent(content)

	w.ShowAndRun()
}
