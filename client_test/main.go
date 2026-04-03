package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

/*
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

		}
	}


		for {
			msg, ok := <-tunnel
			if !ok {
				return
			}
			server_poll(conn, msg)
		}

}
*/

type myTheme struct {
	fyne.Theme
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		return color.RGBA{R: 175, G: 238, B: 238, A: 255} // Красный
	}
	return m.Theme.Color(name, variant)
}

func show_display(msg string, display *widget.Entry) {
	// SEND GROUP
	data := strings.Split(msg, " ")
	if data[0] == "SEND" && data[1] == "GROUP" {
		login := data[2]
		messages := ""
		for i := 4; i < len(data); i++ {
			messages += data[i] + " "
		}
		result := login + ": " + messages
		fyne.Do(func() {
			display.SetText(display.Text + result + "\n")
			display.CursorRow = len(strings.Split(display.Text, "\n")) - 1
			// Принудительно обновляем отображение, чтобы скролл сдвинулся к курсору
			display.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
		})
	}
}

func polling_recv(conn net.Conn, display *widget.Entry) {
	send(conn, "GET GROUP 1")
	//Общий канал получения сообщений
	for {
		msg := recv(conn)
		if msg == "error" {
			println("Ошибка при получении сообщения")
			return
		}
		show_display(msg, display)
	}
}

func login(conn net.Conn, login string, password string, display *widget.Entry) bool {
	send(conn, "login")
	recv(conn)
	send(conn, login)
	recv(conn)
	send(conn, password)
	recv(conn)
	send(conn, "check")
	data := recv(conn)
	if data == "OK" {
		go polling_recv(conn, display)
		return true
	} else {
		return false
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

func key(conn net.Conn, key string) {
	send(conn, key)
	recv(conn)
}

func chat(ip string, logins string, password string) {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{theme.DefaultTheme()})
	myWindow := myApp.NewWindow("Простой Чат")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Область чата (только для чтения)
	chatDisplay := widget.NewMultiLineEntry()
	chatDisplay.Disable()
	chatDisplay.Wrapping = fyne.TextWrapWord

	chatDisplay.SetPlaceHolder("История сообщений...")

	// Поле ввода
	messageInput := widget.NewEntry()
	messageInput.SetPlaceHolder("Введите сообщение...")

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	key(conn, "9S2oPsZJ1ipUxKlbyJvr")
	if !login(conn, logins, password, chatDisplay) {
		println("Не правильный логин или пароль")
		return
	}

	// Функция отправки
	sendMessage := func() {
		text := messageInput.Text
		if text == "" {
			return
		}

		send(conn, "SEND GROUP "+logins+" 1 "+text)
		messageInput.SetText("")
	}

	// Кнопка отправки
	sendButton := widget.NewButton("Отправить", sendMessage)

	// Отправка по Enter
	messageInput.OnSubmitted = func(s string) {
		sendMessage()
	}

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

func main() {
	var ip, login, password string
	fmt.Print("Введите ip: ")
	fmt.Scanln(&ip)
	fmt.Print("Введите логин: ")
	fmt.Scanln(&login)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&password)
	chat(ip, login, password)
}
