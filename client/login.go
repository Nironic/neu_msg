package main

import (
	"github.com/therecipe/qt/widgets"
)

type LoginWindow struct {
	*widgets.QMainWindow
	login    *widgets.QLineEdit
	password *widgets.QLineEdit
	ip       *widgets.QLineEdit
	port     *widgets.QLineEdit
	key      *widgets.QLineEdit
}

func NewLoginWindow() *LoginWindow {
	// Создаем окно
	window := &LoginWindow{
		QMainWindow: widgets.NewQMainWindow(nil, 0),
	}

	window.SetWindowTitle("Вход в чат")
	window.SetMinimumSize2(400, 350)

	// Центральный виджет
	central := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(central)

	// Вертикальное расположение
	layout := widgets.NewQVBoxLayout2(central)
	layout.SetSpacing(10)

	// Заголовок
	title := widgets.NewQLabel2("Настройки подключения", nil, 0)
	title.SetAlignment(1)
	layout.AddWidget(title, 0, 0)

	// Поле IP
	ipLabel := widgets.NewQLabel2("IP сервера:", nil, 0)
	layout.AddWidget(ipLabel, 0, 0)
	window.ip = widgets.NewQLineEdit(nil)
	window.ip.SetText("127.0.0.1")
	layout.AddWidget(window.ip, 0, 0)

	// Поле Port
	portLabel := widgets.NewQLabel2("Порт:", nil, 0)
	layout.AddWidget(portLabel, 0, 0)
	window.port = widgets.NewQLineEdit(nil)
	window.port.SetText("8080")
	layout.AddWidget(window.port, 0, 0)

	// Поле Логин
	loginLabel := widgets.NewQLabel2("Логин:", nil, 0)
	layout.AddWidget(loginLabel, 0, 0)
	window.login = widgets.NewQLineEdit(nil)
	layout.AddWidget(window.login, 0, 0)

	// Поле Пароль
	passLabel := widgets.NewQLabel2("Пароль:", nil, 0)
	layout.AddWidget(passLabel, 0, 0)
	window.password = widgets.NewQLineEdit(nil)
	window.password.SetEchoMode(widgets.QLineEdit__Password)
	layout.AddWidget(window.password, 0, 0)

	// Поле key
	keyLabel := widgets.NewQLabel2("Ключ:", nil, 0)
	layout.AddWidget(keyLabel, 0, 0)
	window.key = widgets.NewQLineEdit(nil)
	layout.AddWidget(window.key, 0, 0)

	layout.AddStretch(0)

	// Кнопка подключения
	connectBtn := widgets.NewQPushButton2("Подключиться", nil)
	connectBtn.SetFixedHeight(35)
	connectBtn.ConnectClicked(func(bool) {
		if window.login.Text() != "" && window.password.Text() != "" {
			// Закрываем окно логина и открываем чат

			ip := window.ip.Text()
			port := window.port.Text()
			login := window.login.Text()
			key := window.key.Text()

			conn, err := connect(ip, port)
			if err == "error" {
				ShowError(window, "Ошибка", "Нет подключения к серверу")
				return
			}
			if !keyserver(conn, key) {
				ShowError(window, "Ошибка", "Не правильный ключ")
				return
			}
			if !loginclient(conn, login, window.password.Text()) {
				ShowError(window, "Ошибка", "Не правильный логин или пароль")
				return
			}

			window.Close()

			// Открываем окно чата
			chatWindow := NewChatWindow(login, conn)
			chatWindow.Show()
		}
	})
	layout.AddWidget(connectBtn, 0, 0)

	return window
}
