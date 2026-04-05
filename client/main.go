package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/therecipe/qt/widgets"
)

func ShowError(parent widgets.QWidget_ITF, title string, message string) {
	msgBox := widgets.NewQMessageBox(parent)
	msgBox.SetWindowTitle(title)
	msgBox.SetText(message)
	msgBox.SetIcon(widgets.QMessageBox__Critical)
	msgBox.SetStandardButtons(widgets.QMessageBox__Ok)
	msgBox.Exec()
}

func main() {
	// Logger configuration
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs // миллисекунды (тоже timestamp)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05", // Формат времени в консоли (только время)
		NoColor:    false,
	}

	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	multi := io.MultiWriter(consoleWriter, file)

	// Настройка логгера с читаемым временем
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	log.Info().Msg("Запуск приложения..")
	app := widgets.NewQApplication(len(os.Args), os.Args)
	SetDarkTheme(app)

	// Создаем и показываем окно логина
	loginWindow := NewLoginWindow()
	loginWindow.Show()

	app.Exec()
}

func SetDarkTheme(app *widgets.QApplication) {
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
		QLineEdit {
			padding: 5px;
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
