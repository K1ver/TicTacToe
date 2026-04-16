# 🎮 TicTacToe (Go + PostgreSQL + Docker)

Полноценное веб-приложение "Крестики-нолики" с backend на Go, REST API, миграциями Goose и SPA фронтендом.

---

## 🚀 Технологии

### Backend
- Go (Golang)
- Chi router
- PostgreSQL
- Goose (миграции)
- Fx (dependency injection)

### Frontend
- Vanilla HTML / CSS / JS (SPA)
- Fetch API

### DevOps
- Docker
- Docker Compose

---

## 📦 Запуск проекта

### 1. Клонировать репозиторий
```bash
git clone git@github.com:K1ver/TicTacToe.git
cd TicTacToe
```

### 2. Создать `.env` по `.env.example`

```env
DATABASE_URL=postgres://postgres:20062009@postgres:5432/tictactoe?sslmode=disable
```

### 3. Запуск через Docker

```bash
docker compose up --build
```

---

## 🌐 Открыть приложение

Frontend:
http://localhost:8080/

API:
http://localhost:8080/api

---

## 🗄️ Миграции

Миграции применяются автоматически через сервис `migrations`.

Полный пересоздание базы:

```bash
docker compose down -v
docker compose up --build
```

---

## 🧩 API

### Auth
- POST /api/auth/signUp
- POST /api/auth/signIn

### User
- GET /api/user/me
- GET /api/user/{userId}

### Game
- POST /api/game
- POST /api/game/{gameId}/join
- POST /api/game/{gameId}
- GET /api/game/{gameId}
- GET /api/game/history
- GET /api/game/browse
- GET /api/game/leaderboard/{count}

---

## 🖥️ Frontend

SPA приложение:

- Авторизация (вход / регистрация)
- Создание игры
- Присоединение к игре
- Игровое поле 3x3
- История игр
- Лидерборд

Файлы:

```
/web/index.html
/web/script.js
/web/style.css
```

---

## 🐘 База данных

Используется PostgreSQL:

- users
- games

Миграции через Goose.

---

## ⚠️ Возможные проблемы

### 400 Bad Request
- неверный JSON
- неправильный UUID
- отсутствует Authorization header

### 404 Not Found
- неправильный endpoint
- неправильный метод (GET/POST)

### Миграции не применились
```bash
docker compose down -v
docker compose up --build
```

---

## 📌 Архитектура

Frontend (SPA)
↓
Go API (Chi router)
↓
Service layer
↓
PostgreSQL

---

## 👨‍💻 Автор

Учебный pet-project для изучения backend разработки:
- Go backend
- REST API
- Docker
- PostgreSQL
- SPA интеграция

---

## 📈 Возможные улучшения

- WebSocket для realtime игры
- React/Vue frontend
- Redis кеширование
- CI/CD pipeline
- Unit tests
- Role-based access control
```