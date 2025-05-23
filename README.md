# 🎓 City Server

**City Server** — это backend-часть уникальной метавселенной, созданной для выпускников школы Летово. Сервер предоставляет WebSocket-коммуникацию, управление игровым миром, а также хостинг Unity WebGL клиента.

---

## 🚀 Возможности

- 🔌 WebSocket сервер для подключения Unity-клиентов
- 🧠 Хаб управления игроками, позициями и взаимодействием
- 🌍 Рассылка `player_joined`, `player_input`, `player_moved` и `player_left`
- 🔐 Авторизация по токену через `Authorization: Bearer`
- 💾 Сохранение информации о токенах и игроках (PostgreSQL)
- 🌐 Хостинг Unity WebGL сборки с поддержкой Brotli-сжатия
- 🧭 Поддержка "наблюдающего сервера" (`PlayerID = SERVER`), который не считается игроком, но получает все сообщения

---

## 📁 Структура проекта

```

cmd/                 # main.go и статическая сборка WebGL
internal/
api/               # REST API (валидация токена и т.п.)
ws/                # WebSocket-сервер, клиент, хаб
services/          # Бизнес-логика
store/             # GORM модели (User, World, Token...)
config/              # YAML конфиги
build/               # Unity WebGL билд

````

---

## ⚙️ Запуск

1. Установи Go ≥ 1.19 и PostgreSQL
2. Создай базу `city_server`
3. Проверь настройки `config/config.yaml`:

```yaml
database:
  dsn: "host=localhost user=postgres password=postgres dbname=city_server sslmode=disable"

server:
  address: "localhost:8080"
  serve_webgl: true
  webgl_path: "./build"
````

4. Собери Unity WebGL и положи его в `./build/`
5. Запусти:

```bash
go run ./cmd
```

---

## 🔌 WebSocket

Подключение:

```
ws://localhost:8080/ws
```

Поддерживаемые сообщения:

* `player_joined`
* `player_input`
* `player_moved`
* `player_left`
* `world_snapshot`

---

## 🛡️ Авторизация

POST `/auth/validate-token`

**Header:**

```
Authorization: Bearer your-token
```

**Response:**

```json
{
  "valid": true,
  "playerId": "player-alex"
}
```

---

## 🧪 Пример fake-клиента

```bash
go run fake-client.go
> input 1 0 true false
> exit
```

---

## 🌐 Хостинг WebGL билда

* Сжатие `.br` поддерживается через `brotliFileServer()`
* Unity билд доступен по адресу:

```
http://localhost:8080/game/
```

---

## 🧠 Поддержка спец-клиента `PlayerID = "SERVER"`

* Может слушать все сообщения
* Не участвует в snapshot или рассылках как обычный игрок

---

## ✨ Проект для

> **City Server** — метавселенная для выпускников Летово. Общение, кампусы университетов, виртуальные встречи, воспоминания и новые связи — в одном пространстве.

---

## 📌 TODO / планы

* Сохранение состояния мира в БД
* Редактор кампусов
* Голосовой чат (WebRTC)
* Расширенная кастомизация игроков
* Мини-игры и квесты

---

## 🧑‍💻 Автор

[Artem Konukhov](https://github.com/IceDarold)
Letovo School ‘26
