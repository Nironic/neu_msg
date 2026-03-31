package main

import (
	"bufio"
	"io"
	"net"
	"strings"

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
		polling(conn, db, user_login)
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

		err := db.CreateUser(user_reg.Login, user_reg.Password, user_reg.Username, "None")
		if err != nil {
			log.Error().Msgf("Ошибка создания пользователя")
			send(conn, "Error")
		}
		send(conn, "OK")

		log.Info().Msgf("Пользователь %s:%s:%s зарегистрирован, адрес: %s", user_reg.Login, user_reg.Password, user_reg.Username, removeAddr)
	}
}

func user_post(conn net.Conn, msg string) { // Основной обработчик сообщений клиента
	if msg == "Hello" {
		send(conn, "Hello the Server")
	}
}

func polling(conn net.Conn, db *MessageDB, user User) {
	defer conn.Close()
	// Общий канал отправки сообщений
	tunnel := make(chan string)
	defer close(tunnel)
	go polling_recv(conn, tunnel)
	for {
		// Ждем сообщения из горутины которая принимает сообщения и обрабатываем их
		msg, ok := <-tunnel
		if !ok {
			return
		}
		user_post(conn, msg)
	}
}

func polling_recv(conn net.Conn, tunnel chan<- string) {
	//Общий канал получения сообщений
	for {
		msg := recv(conn)
		log.Info().Msgf("%s", msg)
		if msg == "error" {
			log.Error().Msgf("Ошибка при получении сообщения [%s]", conn.RemoteAddr().String())
			return
		}
		tunnel <- msg
	}
}
