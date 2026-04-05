package main

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Connection struct {
	login string
	conn  net.Conn
	group int
}

var (
	clients   []*Connection
	clientsMu sync.Mutex
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
		log.Error().Msg("Не могу отправить сообщение клиенту")
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

func loginserver(conn net.Conn, db *MessageDB, rm string, result *User) bool {
	login_client := recv(conn)
	if login_client == "error" {
		log.Error().Msgf("Ошибка при получении логина [%s]", rm)
		return false
	}
	send(conn, "OK")
	password_client := recv(conn)
	if password_client == "error" {
		log.Error().Msgf("Ошибка при получении пароля [%s]", rm)
		return false
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
			return false
		}
		if data.Password == password_client && data.Login == login_client {
			send(conn, "OK") // Пользователь найден и пароль верный
			result.Login = data.Login
			result.Password = data.Password
			result.PathData = data.PathData
			result.Username = data.Username
			// Добавляем клиента
			clientsMu.Lock()
			clients = append(clients, &Connection{
				login: result.Login,
				conn:  conn,
				group: 1,
			})
			clientsMu.Unlock()
			return true
		}
	}
	return false
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
		if loginserver(conn, db, removeAddr, &user_login) {
			log.Info().Msgf("Авторизация [%s]", removeAddr)
			polling(conn, db, user_login)
		}
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

// Протокол взаимодействия
func get_group(conn net.Conn, db *MessageDB, id_group string) {
	idGroup, err := strconv.Atoi(id_group)
	if err != nil {
		log.Error().Msgf("Ошибка при получении id группы [%s]", conn.RemoteAddr().String())
	}
	rows, err := db.db.Query("SELECT id, id_group, user, message, dt, tm FROM groups WHERE id_group = ?", idGroup)
	if err != nil {
		log.Error().Msgf("Ошибка при получении сообщений из группы [%s] [%s]", conn.RemoteAddr().String(), err)
	}
	for rows.Next() {
		var id, idGroup int
		var user, message, dt, tm string
		rows.Scan(&id, &idGroup, &user, &message, &dt, &tm)
		sending := "SEND GROUP " + user + " " + strconv.Itoa(idGroup) + " " + message
		send(conn, sending)
		time.Sleep(10 * time.Millisecond)
	}
}

func user_post(conn net.Conn, msg string, db *MessageDB) { // Основной обработчик сообщений клиента
	// Протокол
	/*
		    Серверу
			Достать данные базы данных личных сообщений
			- GET PERSONAL <login>
			Достать данные из базы данных бесед
			- GET GROUP <id_group>
			Отправить сообщение в личные сообщения
			- SEND PERSONAL <login> <login2> <message>
			Отправить сообщение в группу
			- SEND GROUP <login> <group> <message>

			Клиенту
			Отправить сообщения в группу
			- SEND GROUP <login> <group> <message>

	*/
	if msg == "ping" {
		send(conn, "pong")
		return
	}
	data := strings.Split(msg, " ")
	if len(data) == 3 {
		if data[0] == "GET" && data[1] == "PERSONAL" {
			//Get Messages from personal (send)
		}
		if data[0] == "GET" && data[1] == "GROUP" {
			// Get group messages (send)
			get_group(conn, db, data[2])
		}
	}
	if data[0] == "SEND" && data[1] == "GROUP" {
		login := data[2]
		group, err := strconv.Atoi(data[3])
		if err != nil {
			log.Error().Msgf("Ошибка преобразования в число [%s]", conn.RemoteAddr().String())
		}
		message := ""
		for i := 4; i < len(data); i++ {
			message += data[i] + " "
		}
		err = db.SendMessageGroup(group, login, message)
		if err != nil {
			log.Error().Msgf("Ошибка при сохранения сообщения в базу данных [%s]", conn.RemoteAddr().String())
		}
		for _, c := range clients {
			if c.group == group {
				send(c.conn, "SEND GROUP "+login+" "+strconv.Itoa(group)+" "+message)
			}
		}
	}
}

func polling(conn net.Conn, db *MessageDB, user User) {
	defer func() {
		// Удаляем клиента при выходе
		clientsMu.Lock()
		for i, c := range clients {
			if c.login == user.Login {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		clientsMu.Unlock()
		conn.Close()
	}()
	// Общий канал отправки сообщений
	tunnel := make(chan string)
	defer close(tunnel)
	go polling_recv(conn, tunnel, user)
	for {
		// Ждем сообщения из горутины которая принимает сообщения и обрабатываем их
		msg, ok := <-tunnel
		if !ok {
			return
		}
		user_post(conn, msg, db)
	}
}

func polling_recv(conn net.Conn, tunnel chan<- string, user User) {
	//Общий канал получения сообщений
	for {
		msg := recv(conn)
		log.Info().Msgf("%s", msg)
		if msg == "error" {
			log.Error().Msgf("Ошибка при получении сообщения [%s]", conn.RemoteAddr().String())
			clientsMu.Lock()
			for i, c := range clients {
				if c.login == user.Login {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			clientsMu.Unlock()
			return
		}
		tunnel <- msg
	}
}
