package main

import (
	"io"
	"net"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Структура для хранения всех клиентов
type ClientManager struct {
	// Карта "ID -> соединение"
	clients map[string]net.Conn

	// Канал для добавления новых клиентов
	addCh chan net.Conn

	// Канал для удаления клиентов
	delCh chan net.Conn

	// Мьютекс для защиты доступа к карте clients
	mutex sync.Mutex
}

func inits(key *string) {
	// Прединициализация сервера
	// Загрузка ключа сервера
	data, _ := os.ReadFile("keyserver.key")
	*key = string(data)
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
	log.Info().Msg("Инициализация..")
	key := ""
	inits(&key)
	log.Info().Msgf("Ключ сервера: [%s]", key)
	// Инициализация БД
	log.Info().Msg("Инициализация БД..")
	db, _ := NewMessageDB("./messages.db")
	defer db.Close()

	log.Info().Msg("Тест базы данных..")
	db.CreateTable()

	log.Info().Msg("Запуск сервера..")

	// Server code
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error().Msg("Ошибка запуска сервера")
		return
	} else {
		log.Info().Msg("Сервер запущен! [0.0.0.0:8080]")
	}
	defer listener.Close()

	for {
		log.Info().Msg("Ожидание подключения..")
		// Принимаем соединение
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Msgf("Ошибка при подключении клиента [%s]", conn.RemoteAddr().String())
			continue
		}
		// Обрабатываем соединение в горутине
		go clientHand(conn, key, db)
	}

}
