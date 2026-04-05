# 💬 GoMessenger - Современный мессенджер на Go + Qt

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Qt Version](https://img.shields.io/badge/Qt-5.15-41CD52?style=for-the-badge&logo=qt)](https://www.qt.io)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)

> **Кроссплатформенный мессенджер с красивым темным интерфейсом и реальным временем работы**

![Chat Preview](https://via.placeholder.com/800x400/2b2b2b/ffffff?text=GoMessenger+Chat+Preview)

## ✨ Особенности

- 🎨 **Темная тема** - современный интерфейс без напряжения глаз
- 🔐 **Авторизация** - безопасный вход с логином и паролем
- 💬 **Реальный чат** - мгновенная отправка сообщений
- 📡 **Собственный протокол** - поверх TCP с разделителем `\n`
- 🖥️ **Кроссплатформенность** - Windows, Linux, macOS
- 🔑 **Работа через ключ** - Обмен ключами
- 💣 **Шифрование** - Шифрование чата и базы данных

## 🚀 Быстрый старт

### Требования

- Go 1.20+
- Qt 5.15 (только для разработки)

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/Nironic/neu_msg.git
cd neu_msg

# Установка зависимостей
go mod init neu_msg
go get github.com/therecipe/qt/widgets

# Сборка
go build -ldflags="-H windowsgui" -o messanger.exe