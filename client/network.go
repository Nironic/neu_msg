package main

import (
	"bufio"
	"io"
	"net"
	"strings"

	"github.com/rs/zerolog/log"
)

func connect(ip string, port string) (net.Conn, string) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Info().Msgf("Не могу подключиться к серверу [%s]", ip+":"+port)
		return conn, "error"
	} else {
		log.Info().Msgf("Успешное подключение к серверу")
	}
	return conn, "OK"
}

func recv(conn net.Conn) string {
	reader := bufio.NewReader(conn)

	// Читаем до символа \n
	message, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Error().Err(err).Msg("Не могу принять сообщение от клиента")
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
		log.Error().Msg("Не могу отправить сообщение клиенту")
		return false
	}
	return true
}

func keyserver(conn net.Conn, key string) bool {
	if !send(conn, key) {
		log.Error().Msg("Не могу отправить ключ серверу")
		return false
	} else {
		if recv(conn) == "KEY" {
			log.Info().Msg("Ключи совпали")
			return true
		} else {
			log.Error().Msg("Ключи не совпали")
			return false
		}

	}
}

func loginclient(conn net.Conn, login string, password string) bool {
	rest := false
	if !send(conn, "login") {
		log.Error().Msg("Не могу отправить команду логина серверу")
		return false
	} else {
		data := recv(conn)
		if data == "error" {
			log.Error().Msg("Не могу принять ответ от сервера")
		} else {
			if data == "OK" {
				log.Info().Msg("Перемещаемся в логин")
				rest = true
			}
		}
	}
	if rest {
		if !send(conn, login) {
			log.Error().Msg("Не могу отправить логин серверу")
			return false
		}
		if recv(conn) == "error" {
			log.Error().Msg("Логин не принят")
			return false
		}
		if !send(conn, password) {
			log.Error().Msg("Не могу отправить пароль серверу")
			return false
		}
		if recv(conn) == "error" {
			log.Error().Msg("Пароль не принят")
			return false
		}
		if !send(conn, "check") {
			log.Error().Msg("Не могу отправить команду check серверу")
			return false
		}
		data := recv(conn)
		if data == "error" {
			log.Error().Msg("Не могу принять ответ от сервера")
			return false
		}
		if data == "OK" {
			log.Info().Msg("Авторизация прошла успешно")
		}
		if data == "NOUSER" {
			log.Error().Msg("Пользователь не найден")
			return false
		}
	}
	return true
}

func polling(conn net.Conn, tunnel chan<- string) bool {
	for {
		msg := recv(conn)
		log.Info().Msgf("%s", msg)
		if msg == "error" {
			log.Error().Msgf("Ошибка при получении сообщения [%s]", conn.RemoteAddr().String())
			return false
		}
		tunnel <- msg
	}
}
