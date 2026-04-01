// database.go
package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Структура для БД
type MessageDB struct {
	db *sql.DB
}

// Message структура
type Message struct {
	ID    int
	User1 string
	User2 string
	Msg   string
	Dt    string
	Tm    string
}

// User структура
type User struct {
	ID       int
	Login    string
	Password string
	Username string
	PathData string
}

// Структура для групп
type Group struct {
	ID       int
	Id_group string
	User     string
	Message  string
	Dt       string
	Tm       string
}

// Создание новой БД
func NewMessageDB(path string) (*MessageDB, error) {
	// Проверяем файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Error().Msgf("Файл БД не найден: %s", path)
		log.Info().Msgf("Создаем файл БД: %s", path)
	}

	// Подключаемся
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MessageDB{db: db}, nil
}

func (m *MessageDB) CreateUser(login string, password string, username string, path_data string) error {
	_, err := m.db.Exec(`
        INSERT INTO users (login, password, username, path_data) 
        VALUES (?, ?, ?, ?)
    `, login, password, username, path_data)
	return err
}

// Закрытие БД
func (m *MessageDB) Close() error {
	return m.db.Close()
}

// Проверка существования таблицы
func (m *MessageDB) TableExists(tableName string) bool {
	var name string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
	err := m.db.QueryRow(query, tableName).Scan(&name)
	return err == nil
}

func (m *MessageDB) CreateTable() error {
	// Проверка таблиц
	if !m.TableExists("message") {
		log.Error().Msg("Таблица 'message' не найдена")
		log.Info().Msg("Создаем таблицу 'message'")
		_, err := m.db.Exec(`
        CREATE TABLE message (
            id    INTEGER PRIMARY KEY AUTOINCREMENT,
            user1 TEXT NOT NULL,
            user2 TEXT NOT NULL,
            msg   TEXT NOT NULL,
            dt    TEXT NOT NULL,
			tm    TEXT NOT NULL
        	)
    	`)
		if err != nil {
			log.Error().Msgf("Ошибка создания таблицы 'message': %s", err)
			return err
		}
		log.Info().Msg("Таблица 'message' создана")
	} else {
		log.Info().Msg("Таблица 'message' найдена")
	}
	if !m.TableExists("users") {
		log.Error().Msg("Таблица 'users' не найдена")
		log.Info().Msg("Создаем таблицу 'users'")
		_, err := m.db.Exec(`
        CREATE TABLE users (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			login     TEXT    NOT NULL,
			password  TEXT    NOT NULL,
			username TEXT    NOT NULL,
			path_data TEXT    NOT NULL
		);`)
		if err != nil {
			log.Error().Msgf("Ошибка создания таблицы 'users': %s", err)
			return err
		}
		log.Info().Msg("Таблица 'users' создана")
	} else {
		log.Info().Msg("Таблица 'users' найдена")
	}
	return nil
}

// Проверка существования таблицы

// Сохранение сообщения
func (m *MessageDB) SendMessage(user1 string, user2 string, msg string) error {
	now := time.Now()
	dt := now.Format("2006-01-02")
	tm := now.Format("15:04:05")
	_, err := m.db.Exec(`
        INSERT INTO message (user1, user2, msg, dt, tm) 
        VALUES (?, ?, ?, ?, ?)
    `, user1, user2, msg, dt, tm)
	return err
}

func (m *MessageDB) SendMessageGroup(id_group string, user string, msg string) error {
	id_group_restruck, err := strconv.Atoi(id_group)
	if err != nil {
		log.Error().Msgf("Ошибка конвертации id_group: %s", err)
	}
	now := time.Now()
	dt := now.Format("2006-01-02")
	tm := now.Format("15:04:05")
	_, err = m.db.Exec(`
        INSERT INTO groups (id_group, user, message, dt, tm) 
        VALUES (?, ?, ?, ?, ?)
    `, id_group_restruck, user, msg, dt, tm)
	return err
}

// Получить данные пользователя
func (m *MessageDB) GetUserAll(login string) (*User, error) {
	var user User
	query := `SELECT id, login, password, username, path_data FROM users WHERE login = ?`
	err := m.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Username, &user.PathData)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
