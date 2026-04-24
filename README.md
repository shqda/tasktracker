# TaskTracker

![CI](https://github.com/shqda/tasktracker/actions/workflows/ci.yml/badge.svg)
![Go version](https://img.shields.io/badge/go-1.25-blue)
![License](https://img.shields.io/badge/license-MIT-green)

> Учебный проект в процессе разработки.

REST API для управления задачами, написанный на Go.

## Стек

- **Go** + **Gin** — HTTP-сервер
- **PostgreSQL** + **sqlx** — хранилище
- **goose** — миграции
- **slog** — логирование
- **Docker** + **docker-compose** — запуск окружения

## Быстрый старт

### Через Docker

```bash
cp .env.example .env
docker compose up
```

Сервис поднимется на `http://localhost:8080`.  
Миграции применяются автоматически при старте.

### Локально

**Требования:** Go 1.25+, PostgreSQL, goose

```bash
# Настройка переменных окружения
cp .env.example .env

# Применить миграции
make migration-up

# Запустить сервер
go run ./cmd/task_tracker
```

## API

Swagger UI доступен по адресу: `http://localhost:8080/swagger/index.html`

| Метод    | Путь          | Описание                  |
|----------|---------------|---------------------------|
| `GET`    | `/tasks `     | Получить все задачи       |
| `GET`    | `/tasks/last` | Получить последнюю задачу |
| `GET`    | `/tasks/{id}` | Получить задачу по ID     |
| `POST`   | `/tasks`      | Создать задачу            |
| `PUT`    | `/tasks/{id}` | Переименовать задачу      |
| `DELETE` | `/tasks/{id}` | Удалить задачу            |

### Примеры

**Создать задачу**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"task": "buy milk"}'
```
```json
{"id": 1, "title": "buy milk"}
```

**Получить все задачи**
```bash
curl http://localhost:8080/tasks 
```
```json
[{"id": 1, "title": "buy milk"}]
```

**Переименовать задачу**
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "buy bread"}'
```

**Удалить задачу**
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Makefile

| Команда                          | Описание                        |
|----------------------------------|---------------------------------|
| `make test`                      | Все тесты                       |
| `make utest`                     | Юнит-тесты                      |
| `make itest`                     | Интеграционные тесты            |
| `make gen`                       | Сгенерировать моки              |
| `make migration-up`              | Применить миграции              |
| `make migration-down`            | Откатить миграции               |
| `make migration-status`          | Статус миграций                 |
| `make migration-create NAME=...` | Создать новый файл миграции     |
