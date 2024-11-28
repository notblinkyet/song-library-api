# Song Library API

Song Library API — это простое RESTful API для управления музыкальными песнями. API поддерживает операции CRUD (создание, чтение, обновление, удаление) и предоставляет возможность фильтрации песен.

## Содержание

- [Описание](#описание)
- [Структура проекта](#структура-проекта)
- [Установка и настройка](#установка-и-настройка)
- [Запуск](#запуск)
- [API Эндпоинты](#api-эндпоинты)
  - [GET /songs](#get-songs)
  - [POST /songs](#post-songs)
  - [GET /songs/{id}](#get-songsid)
  - [PATCH /songs/{id}](#patch-songsid)
  - [DELETE /songs/{id}](#delete-songsid)
- [Известные проблемы](#известные-проблемы)
- [Технологии](#технологии)

---

## Описание

Song Library API позволяет:

- Создавать новые записи о песнях.
- Просматривать список песен с фильтрацией по различным параметрам.
- Получать подробную информацию о конкретной песне.
- Обновлять и удалять записи.

API разработано с использованием Go (Golang), а база данных — PostgreSQL. Также интегрировано стороннее API для получения дополнительной информации о песнях.

---

## Структура проекта
├── cmd
│   ├── app
│   │   └── main.go         # Точка входа для запуска HTTP-сервера
│   └── migrator
│       └── main.go         # Точка входа для запуска миграций
├── docs                    # Swagger-документация
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── config              # Загрузка конфигураций из .env
│   │   └── config.go
│   ├── database            # Работа с базой данных
│   │   ├── database.go     # Интерфейс базы данных
│   │   ├── migrations      # SQL-миграции
│   │   └── postgresql      # Реализация работы с PostgreSQL
│   ├── lib
│   │   ├── api             # Клиент для стороннего API
│   │   ├── ParseURL        # Парсинг URL параметров
│   │   └── sl              # Логирование ошибок
│   ├── logger              # Конфигурация логгера
│   ├── models              # Определение моделей данных
│   ├── services            # Основная бизнес-логика
│   └── transport
│       └── http            # HTTP-хэндлеры и эндпоинты

---

## Установка и настройка

1. Клонируйте репозиторий:

     git clone https://github.com/notblinkyet/song-library-api.git
   cd song-library-api
   
2. Установите зависимости:

   Убедитесь, что Go установлен, и выполните команду:

     go mod tidy
   
3. Настройте `.env` файл:

   Создайте .env файл в корневой директории и заполните его:

   DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=song_library
    DB_USER=postgres
    DB_PASSWORD=12345pass
    MIGRATION_PATH=/home/hobonail/go_projects/song-library-api/internal/database/migrations
    SERVER_PORT=9090
    SERVER_HOST=localhost
    TIMEOUT=5
    IDLE_TIMEOUT=30
    API_ADDR_URL=http://example_api
   
---

## Запуск

### Запуск HTTP-сервера

Для запуска API выполните команду:
go run cmd/app/main.go

Сервер будет доступен по адресу: http://localhost:9090.

### Применение миграций

Для применения миграций выполните:
go run cmd/migrator/main.go

---

Так же есть параметр --rollback - количество миграций для отката

## API Эндпоинты

### GET /songs

Получение списка песен с возможностью фильтрации.

#### Пример запроса через curl:
curl -X 'GET' \
  'http://localhost:9090/songs?song=Supermassive%20Black%20Hole&group=Muse&release_date=16.07.2006&link=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DXsp3_a-PMTw' \
  -H 'accept: application/json'

---

### POST /songs

Создание новой песни.

#### Пример запроса через curl:
curl -X 'POST' \
  'http://localhost:9090/songs' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}'

---

### GET /songs/{id}

Получение информации о песне по ID.

#### Пример запроса через curl:
curl -X 'GET' \
  'http://localhost:9090/songs/13?start=1&count=1' \
  -H 'accept: application/json'

---

### PATCH /songs/{id}

Обновление информации о песне по ID.

#### Пример запроса через curl:
curl -X 'PATCH' \
  'http://localhost:9090/songs/11' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "group": "Muse",
  "id": 0,
  "link": "https://example.com",
  "releaseDate": "2006-07-16",
  "song": "Supermassive Black Hole",
  "text": "Some lyrics"
}'

---

### DELETE /songs/{id}

Удаление песни по ID.

#### Пример запроса через curl:
curl -X 'DELETE' \
  'http://localhost:9090/songs/10' \
  -H 'accept: application/json'

---

## Известные проблемы

- Фильтр по тексту (параметр `text`) на эндпоинте GET `/songs` работает некорректно. Требуется доработка для улучшения функциональности.

---

## Технологии

- Golang — язык программирования.
- PostgreSQL — реляционная база данных.
- Chi — маршрутизация.
- Swaggo — генерация Swagger-документации.
- pgxpool — пул соединений для PostgreSQL.