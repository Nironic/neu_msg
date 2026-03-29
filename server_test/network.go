package main

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

func recv(conn net.Conn) string {
	reader := bufio.NewReader(conn)

	// Читаем до символа \n
	message, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Error().Err(err).Msg("Ошибка чтения")
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
		log.Error().Msg("Не могу принять сообщение от клиента")
		return false
	}
	return true
}

// Key
func keyserver(conn net.Conn, key string, raddr string) bool {
	key_client := recv(conn) // Отправляет ключ
	if key_client == "error" {
		log.Error().Msgf("Ошибка при получении ключа [%s]", raddr)
		return false
	}
	if key_client == key {
		send(conn, "KEY") // Ключи совпали
		log.Info().Msgf("Ключи совпали [%s]", raddr)
		return true
	} else {
		send(conn, "NOKEY") // Ключи не совпали
		log.Error().Msgf("Ключи не совпали [%s]", raddr)
		return false
	}
}

func loginserver(conn net.Conn, db *MessageDB, rm string, result *User) {
	login_client := recv(conn)
	if login_client == "error" {
		log.Error().Msgf("Ошибка при получении логина [%s]", rm)
		return
	}
	send(conn, "OK")
	password_client := recv(conn)
	if password_client == "error" {
		log.Error().Msgf("Ошибка при получении пароля [%s]", rm)
		return
	}
	send(conn, "OK")
	check := recv(conn)
	if check == "error" {
		log.Error().Msgf("Ошибка при получении ответа [%s]", rm)
	}
	if check == "check" {
		data, err := db.GetUserAll(login_client)
		if err != nil {
			log.Error().Msgf("Пользователя %s не существует [%s] %s", login_client, rm, err)
			send(conn, "NOUSER")
			return
		}
		if data.Password == password_client && data.Login == login_client {
			send(conn, "OK") // Пользователь найден и пароль верный
			result.Login = data.Login
			result.Password = data.Password
			result.PathData = data.PathData
			result.Username = data.Username
		}
	}

}

// Hadnler Client
func clientHand(conn net.Conn, key string, db *MessageDB) {
	defer conn.Close() // Предварительное закрытие соединения для оптимизации
	defer log.Info().Msgf("Отключен [%s]", conn.RemoteAddr().String())
	removeAddr := conn.RemoteAddr().String()
	log.Info().Msgf("Подключен [%s]", removeAddr)
	// Обмен ключами
	if !keyserver(conn, key, removeAddr) {
		return
	}
	// Принимаем либо логин либо регистрацию
	command := recv(conn)
	if command == "error" {
		log.Error().Msgf("Ошибка при получении команды [%s]", removeAddr)
		return
	}
	if command == "login" {
		user_login := User{}
		send(conn, "OK")
		loginserver(conn, db, removeAddr, &user_login)
		log.Info().Msgf("Авторизация [%s]", removeAddr)
		polling(conn, db)
	}
	if command == "reg" {
		log.Info().Msgf("Регистрация [%s]", removeAddr)
		send(conn, "OK")
		user_reg := User{}
		user_reg.Login = recv(conn)
		if user_reg.Login == "error" {
			log.Error().Msgf("Ошибка при получении логина [%s]", removeAddr)
		}
		send(conn, "OK")
		user_reg.Password = recv(conn)
		if user_reg.Password == "error" {
			log.Error().Msgf("Ошибка при получении пароля [%s]", removeAddr)
		}
		send(conn, "OK")
		user_reg.Username = recv(conn)
		if user_reg.Username == "error" {
			log.Error().Msgf("Ошибка при получении имени пользователя [%s]", removeAddr)
		}
		send(conn, "OK")

		db.CreateUser(user_reg.Login, user_reg.Password, user_reg.Username, "None")
		log.Info().Msgf("Пользователь %s:%s:%s зарегистрирован, адрес: %s", user_reg.Login, user_reg.Password, user_reg.Username, removeAddr)
	}
}

func polling(conn net.Conn, db *MessageDB) {
	var wg sync.WaitGroup
	wg.Add(1)
	tunnel := make(chan string)
	run := true
	go func() {
		defer wg.Done()
		client_long_pooling(conn, db, &run, tunnel)
	}()
	client_short_pooling(conn, db, &run, tunnel)
	wg.Wait()
}

// Два основных потока
// Постоянно принимает сообщения
func client_long_pooling(conn net.Conn, db *MessageDB, run *bool, tunnel chan<- string) {
	for {
		if !*run {
			log.Info().Msgf("Выход из горутины чтения.")
			close(tunnel) // Убеждаемся, что канал закрыт при выходе
			return
		}
		data := recv(conn)
		if data == "error" {
			*run = false
			log.Error().Msgf("Ошибка при получении сообщения [%s]", conn.RemoteAddr().String())
			close(tunnel)
		}

		tunnel <- data
	}
}

// Постоянно отправляет сообщения
func client_short_pooling(conn net.Conn, db *MessageDB, run *bool, tunnel <-chan string) {
	for *run {
		revers := <-tunnel

		if revers == "syns" {
			if !send(conn, "OK") {
				log.Error().Msgf("Ошибка ответа синхронизации %s", conn.RemoteAddr().String())
			}
		}
	}
}
