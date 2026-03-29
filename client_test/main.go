package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
