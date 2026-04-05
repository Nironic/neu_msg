package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/therecipe/qt/widgets"
)

type ChatWindow struct {
	*widgets.QMainWindow
	messages *widgets.QTextEdit
	input    *widgets.QLineEdit
}

func NewChatWindow(login string, conn net.Conn) *ChatWindow {
	// Создаем окно
	window := &ChatWindow{
		QMainWindow: widgets.NewQMainWindow(nil, 0),
	}

	window.SetWindowTitle(fmt.Sprintf("Чат - %s", login))
	window.SetMinimumSize2(600, 400)

	// Центральный виджет
	central := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(central)

	// Вертикальное расположение
	layout := widgets.NewQVBoxLayout2(central)

	// Поле для сообщений (только чтение)
	window.messages = widgets.NewQTextEdit(nil)
	window.messages.SetReadOnly(true)
	layout.AddWidget(window.messages, 0, 0)
	window.messages.Font().SetPointSize(12)

	// Получение сообщений
	tunnel := make(chan string)
	go polling(conn, tunnel)
	go ShowMessages(window, tunnel)
	send(conn, "CONNECT 1"+login)
	send(conn, "GET GROUP 1")

	// Горизонтальная панель для ввода
	inputLayout := widgets.NewQHBoxLayout2(nil)

	window.input = widgets.NewQLineEdit(nil)
	window.input.SetPlaceholderText("Введите сообщение...")

	sendBtn := widgets.NewQPushButton2("Отправить", nil)

	inputLayout.AddWidget(window.input, 0, 0)
	inputLayout.AddWidget(sendBtn, 0, 0)
	layout.AddLayout(inputLayout, 0)

	// Отправка по кнопке
	sendBtn.ConnectClicked(func(bool) {
		window.SendMessage(conn, login)
	})

	// Отправка по Enter
	window.input.ConnectReturnPressed(func() {
		window.SendMessage(conn, login)
	})

	return window
}

func (w *ChatWindow) SendMessage(conn net.Conn, login string) {
	text := w.input.Text()
	if text != "" {
		result := "SEND GROUP " + login + " 1 " + text
		send(conn, result)
		w.input.Clear()
	}
}

func protocol(data string) (string, string) {
	msg := strings.Split(data, " ")
	if msg[0] == "SEND" && msg[1] == "GROUP" {
		login := msg[2]
		message := ""
		for i := 4; i < len(msg); i++ {
			message += msg[i] + " "
		}
		return login, message
	}
	return "None", "None"
}

func ShowMessages(w *ChatWindow, tunnel chan string) {
	for {
		msg, ok := <-tunnel
		if !ok {
			ShowError(w, "Ошибка", "Подключение с сервером утеряно")
			break
		}
		if msg != "" {
			login, message := protocol(msg)
			w.messages.Append(fmt.Sprintf("%s: %s", login, message))
		}
	}
}
