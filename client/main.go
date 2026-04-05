package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

type LoginWindow struct {
	*widgets.QMainWindow
	login    *widgets.QLineEdit
	password *widgets.QLineEdit
	ip       *widgets.QLineEdit
	port     *widgets.QLineEdit
}

type ChatWindow struct {
	*widgets.QMainWindow
	messages *widgets.QTextEdit
	input    *widgets.QLineEdit
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Включаем темную тему через CSS
	SetDarkTheme(app)

	loginWindow := NewLoginWindow()
	loginWindow.Show()

	app.Exec()
}

func SetDarkTheme(app *widgets.QApplication) {
	// Простая темная тема через CSS
	app.SetStyleSheet(`
		QMainWindow, QWidget {
			background-color: #2b2b2b;
			color: #ffffff;
		}
		QPushButton {
			background-color: #3c3c3c;
			border: 1px solid #555;
			padding: 8px;
			border-radius: 4px;
			color: white;
			font-weight: bold;
		}
		QPushButton:hover {
			background-color: #4a4a4a;
		}
		QPushButton:pressed {
			background-color: #2a2a2a;
		}
		QLineEdit {
			padding: 5px;
			border: 1px solid #555;
			border-radius: 3px;
			background-color: #3c3c3c;
			color: white;
		}
		QLineEdit:focus {
			border: 1px solid #5a8aba;
		}
		QTextEdit {
			border: 1px solid #555;
			border-radius: 3px;
			background-color: #3c3c3c;
			color: white;
		}
		QLabel {
			color: #ffffff;
		}
	`)
}

func NewLoginWindow() *LoginWindow {
	window := &LoginWindow{
		QMainWindow: widgets.NewQMainWindow(nil, 0),
	}

	window.SetWindowTitle("Вход в чат")
	window.SetMinimumSize2(400, 350)

	central := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(central)

	layout := widgets.NewQVBoxLayout2(central)
	layout.SetSpacing(10)

	// Заголовок
	title := widgets.NewQLabel2("Настройки подключения", nil, 0)
	title.SetAlignment(1)
	title.SetStyleSheet("font-size: 16px; font-weight: bold; margin-bottom: 10px;")
	layout.AddWidget(title, 0, 0)

	// IP
	ipLabel := widgets.NewQLabel2("IP сервера:", nil, 0)
	layout.AddWidget(ipLabel, 0, 0)
	window.ip = widgets.NewQLineEdit(nil)
	window.ip.SetText("127.0.0.1")
	layout.AddWidget(window.ip, 0, 0)

	// Port
	portLabel := widgets.NewQLabel2("Порт:", nil, 0)
	layout.AddWidget(portLabel, 0, 0)
	window.port = widgets.NewQLineEdit(nil)
	window.port.SetText("8080")
	layout.AddWidget(window.port, 0, 0)

	// Login
	loginLabel := widgets.NewQLabel2("Логин:", nil, 0)
	layout.AddWidget(loginLabel, 0, 0)
	window.login = widgets.NewQLineEdit(nil)
	layout.AddWidget(window.login, 0, 0)

	// Password
	passLabel := widgets.NewQLabel2("Пароль:", nil, 0)
	layout.AddWidget(passLabel, 0, 0)
	window.password = widgets.NewQLineEdit(nil)
	window.password.SetEchoMode(widgets.QLineEdit__Password)
	layout.AddWidget(window.password, 0, 0)

	layout.AddStretch(0)

	// Кнопка подключения
	connectBtn := widgets.NewQPushButton2("Подключиться", nil)
	connectBtn.SetFixedHeight(35)
	connectBtn.ConnectClicked(func(bool) {
		if window.login.Text() != "" && window.password.Text() != "" {
			window.ConnectToServer()
		}
	})
	layout.AddWidget(connectBtn, 0, 0)

	return window
}

func (w *LoginWindow) ConnectToServer() {
	ip := w.ip.Text()
	port := w.port.Text()
	login := w.login.Text()
	password := w.password.Text()
	println(password)

	fmt.Printf("Подключение к %s:%s как %s\n", ip, port, login)

	w.Close()

	chatWindow := NewChatWindow(login)
	chatWindow.Show()
}

func NewChatWindow(username string) *ChatWindow {
	window := &ChatWindow{
		QMainWindow: widgets.NewQMainWindow(nil, 0),
	}

	window.SetWindowTitle(fmt.Sprintf("Чат - %s", username))
	window.SetMinimumSize2(600, 400)

	central := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(central)

	layout := widgets.NewQVBoxLayout2(central)

	// Поле сообщений
	window.messages = widgets.NewQTextEdit(nil)
	window.messages.SetReadOnly(true)
	layout.AddWidget(window.messages, 0, 0)

	// Панель ввода
	inputLayout := widgets.NewQHBoxLayout2(nil)
	window.input = widgets.NewQLineEdit(nil)
	window.input.SetPlaceholderText("Введите сообщение...")
	sendBtn := widgets.NewQPushButton2("Отправить", nil)

	inputLayout.AddWidget(window.input, 0, 0)
	inputLayout.AddWidget(sendBtn, 0, 0)
	layout.AddLayout(inputLayout, 0)

	sendBtn.ConnectClicked(func(bool) {
		window.SendMessage()
	})

	window.input.ConnectReturnPressed(func() {
		window.SendMessage()
	})

	return window
}

func (w *ChatWindow) SendMessage() {
	text := w.input.Text()
	if text != "" {
		w.messages.Append("Я: " + text)
		w.input.Clear()
	}
}
