
# Song Library API

Простое RESTful API для управления музыкальными песнями.

## Содержание

1. [Описание](#описание)
2. [Структура проекта](#структура-проекта)
3. [Запуск проекта](#запуск-проекта)
4. [Примеры запросов](#примеры-запросов)
5. [Заметки](#заметки)

---

## Описание

Song Library API предоставляет возможность управлять библиотекой песен, включая создание, чтение, обновление и удаление данных о песнях.

---

## Структура проекта

```
├── cmd
│   ├── app
│   │   └── main.go        # Точка входа, где поднимается HTTP сервер
│   └── migrator
│       └── main.go        # Точка входа для применения миграций
├── docs
│   ├── docs.go            # Сгенерировано утилитой swag
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod                 # Зависимости проекта
├── go.sum
├── internal
│   ├── config
│   │   └── config.go      # Структура конфигурации и загрузка из .env
│   ├── database
│   │   ├── database.go    # Интерфейс базы данных
│   │   ├── migrations     # SQL файлы миграций
│   │   └── postgresql     # Реализация работы с PostgreSQL
│   ├── lib
│   │   ├── api
│   │   │   └── api.go     # Клиент для работы с внешним API
│   │   ├── ParseURL       # Парсинг значений из URL
│   │   └── sl             # Логгер ошибок
│   ├── logger
│   │   └── logger.go      # Настройка логгера
│   ├── models
│   │   └── models.go      # Структуры данных для базы и запросов
│   ├── services
│   │   └── service.go     # Бизнес-логика
│   └── transport
│       └── http           # HTTP хендлеры и эндпоинты
└── README.md              # Документация
```

---

## Запуск проекта

1. Установите зависимости:
   ```bash
   go mod tidy
   ```

2. Запустите сервер:
   ```bash
   go run cmd/app/main.go
   ```

3. Примените миграции:
   ```bash
   go run cmd/migrator/main.go
   ```

### Пример `.env` файла

```env
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
```

### Откат миграций

Для отката миграций используйте параметр `--rollback` с указанием количества миграций для отката:
```bash
go run cmd/migrator/main.go --rollback 1
```

---

## Примеры запросов

### Создание песни

**POST** `/songs`

Пример запроса через `curl`:

```bash
curl -X 'POST'   'http://localhost:9090/songs'   -H 'accept: application/json'   -H 'Content-Type: application/json'   -d '{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}'
```

---

### Получение песен с фильтром

**GET** `/songs`

Пример запроса через `curl`:

```bash
curl -X 'GET'   'http://localhost:9090/songs?song=Supermassive%20Black%20Hole&group=Muse'   -H 'accept: application/json'
```

---

### Получение текста песни

**GET** `/songs/{id}`

Пример запроса через `curl`:

```bash
curl -X 'GET'   'http://localhost:9090/songs/13?start=1&count=1'   -H 'accept: application/json'
```

---

### Удаление песни

**DELETE** `/songs/{id}`

Пример запроса через `curl`:

```bash
curl -X 'DELETE'   'http://localhost:9090/songs/10'   -H 'accept: application/json'
```

---

### Обновление песни

**PATCH** `/songs/{id}`

Пример запроса через `curl`:

```bash
curl -X 'PATCH'   'http://localhost:9090/songs/11'   -H 'accept: application/json'   -H 'Content-Type: application/json'   -d '{
  "group": "Muse",
  "id": 0,
  "link": "string",
  "releaseDate": "string",
  "song": "string",
  "text": "string"
}'
```

---

## Заметки

- Фильтр по тексту в `GET` запросе `/songs` работает не всегда корректно и требует доработки.
- Документация Swagger доступна по адресу: [http://localhost:9090/swagger/index.html](http://localhost:9090/swagger/index.html).