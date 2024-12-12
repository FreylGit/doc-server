# Doc Server

## Описание
`Doc Server` — это сервер для управления документами с использованием REST API. Поддерживает аутентификацию с токенами JWT и предоставляет API для создания, получения, редактирования и удаления документов.

---

## Запуск приложения

Перед запуском убедитесь, что файл конфигурации `local.env` корректно настроен и миграции выполнены.

### Команда запуска
```bash
# bash
docker compose --env-file local.env up -d
```
Так же в репозитории есть файл для импорта в postman `docs_server.postman_collection.json`
## Конфигурация

Пример файла `local.env`:
```
HTTP_PORT=8081
HTTP_HOST=localhost

PG_DATABASE_NAME=doc
PG_USER=doc-user
PG_PASSWORD=doc-password
PG_HOST=localhost
PG_PORT=5435
MIGRATION_DIR=./migrations
PG_DSN="host=localhost port=5435 dbname=doc user=doc-user password=doc-password sslmode=disable"
DB_HOST_CONTAINER=pg_doc
PG_PORT_CONTAINER=5432

PG_MAX_CONNS=20
PG_MIN_CONNS=5
PG_MAX_CONN_LIFETIME=5
PG_HEALTH_CHECK_PERIOD=1

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=12345
REDIS_CACHE_TTL=3600

APP_TOKEN_ADMIN=fewfw142veKfwq24
APP_SECRET_KEY=doc_secret_key
```

---

## Зависимости

В проекте используются следующие библиотеки (указаны в `go.mod`):
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- [jwt-go](https://github.com/dgrijalva/jwt-go) — для работы с токенами JWT
- [PostgreSQL](https://github.com/jackc/pgx) — драйвер для работы с базой данных
- [Redis](https://github.com/redis/go-redis) — клиент для Redis
- [Zap](https://go.uber.org/zap) — для логирования

Полный список зависимостей см. в `go.mod`.

---

### Ручки API

#### 1. Авторизация
```http
POST http://localhost:8081/api/auth
Content-Type: application/json

{
    "login": "test@mail.ru",
    "pswd": "1FEwwefwek4Evr@"
}
```

#### 2. Регистрация
```http
POST http://localhost:8081/api/register
Content-Type: application/json

{
    "token": "fewfw142veKfwq24",
    "login": "test@mail.ru",
    "pswd": "1FEwwefwek4Evr@"
}
```

#### 3. Добавление документа
```http
POST http://localhost:8081/api/docs
Content-Type: multipart/form-data

meta={
    "name": "sample1.json",
    "file": true,
    "public": false,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQxODM5OTgsImp0aSI6IjEiLCJzdWIiOiJ0ZXN0QG1haWwucnUifQ.9hdzN-f9Tr16AJ2vx9a6doU18o9B5TWkr9ZBPq6qNls",
    "mime": "json",
    "grant": ["test@mail.ru"]
}
file=@path/to/your/file.json
```

#### 4. Получение документа
```http
GET http://localhost:8081/api/docs/:id?token= <your_jwt_token>
```
#### 5. Получение докуменов
```http
GET http://localhost:8081/api/docs/?token= <your_jwt_token>&limit=<limit>&....
```
#### 6. Удаление документа
```http
DELETE http://localhost:8081/api/docs/:id??token= <your_jwt_token>
```

---

## Миграции

Миграции хранятся в директории `migrations`. Пример начальной миграции:

---

## Разработка

### Структура проекта
```
.
├── cmd
│   └── main.go
├── internal
│   ├── api
│   │   ├── handlers
│   │   ├── routes.go
│   ├── config
│   ├── models
│   ├── services
│   ├── storage
│   └── utils
├── migrations
├── local.env
├── docker-compose.yaml
└── docs_server.postman_collection.json
```

