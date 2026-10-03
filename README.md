# URL Shortener (Go)

REST API для сокращения ссылок на Go, реализованный на стандартной библиотеке `net/http`, без сторонних фреймворков. Хранилище спрятано за интерфейсом — поддерживаются in-memory и PostgreSQL реализации, переключаются одной строкой в `main.go`.

## Возможности

- Сокращение ссылки с генерацией уникального кода (6 символов, буквы+цифры)
- Редирект по короткой ссылке с подсчётом переходов
- Получение информации о ссылке без перехода
- Обновление оригинального URL
- Удаление короткой ссылки
- Статистика переходов (количество и время последнего доступа)
- Валидация входного URL (схема http/https, наличие хоста)

## Технологии

- Go, стандартная библиотека: `net/http`, `encoding/json`, `net/url`, `crypto/rand`, `database/sql`
- Драйвер PostgreSQL: `jackc/pgx/v5/stdlib`
- Переменные окружения из `.env`: `joho/godotenv`
- Роутинг на чистом `http.HandleFunc` (без сторонних роутеров), с ручным разбором метода и параметров пути
- Хранилище спрятано за интерфейсом `Storage`: `MemoryStorage` (map в памяти) и `PostgresStorage` (SQL, параметризованные запросы)

## Архитектура

```
url-shortener/
├── main.go
├── init.sql                      # схема таблицы urls
├── .env.example                  # шаблон переменных окружения
├── models/
│   └── url.go                    # структура URLRecord
├── internal/
│   ├── storage/
│   │   ├── storage.go            # интерфейс Storage
│   │   ├── memory_storage.go     # реализация на map
│   │   ├── postgres_storage.go   # реализация на PostgreSQL
│   │   └── id_generator.go       # генерация случайного кода
│   └── handlers/
│       └── handlers.go           # HTTP-обработчики
└── dto/
    └── Dtos.go                   # структуры запросов/ответов
```

## Запуск

### 1. База данных

```bash
psql -h localhost -p 5432 -U <username>
CREATE DATABASE urlshortener;
\c urlshortener
\i init.sql
```

### 2. Переменные окружения

Скопируй `.env.example` в `.env` и заполни реальными значениями:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=<username>
DB_PASSWORD=<password>
DB_NAME=urlshortener
```

### 3. Запуск сервера

```bash
git clone https://github.com/flecc1/URL-Shortener.git
cd URL-Shortener
go run main.go
```
Сервер поднимается на `http://localhost:8080`.

## Эндпоинты

| Метод  | Путь                  | Описание                         |
|--------|-----------------------|----------------------------------|
| POST   | `/shorten`            | Создать короткую ссылку          |
| GET    | `/{id}`               | Перейти по короткой ссылке (302) |
| GET    | `/shorten/{id}`       | Получить информацию о ссылке     |
| PUT    | `/shorten/{id}`       | Обновить оригинальный URL        |
| DELETE | `/shorten/{id}`       | Удалить ссылку                   |
| GET    | `/shorten/{id}/stats` | Статистика переходов             |
| GET    | `/all`                | Вытащить все ссылки              |

## Пример использования

**Создание ссылки:**
```bash
curl -X POST http://localhost:8080/shorten \
  -d '{"url": "https://google.com"}'
```
Ответ (201):
```json
{
	"id": "AULm6Z",
	"short_url": "http://localhost:8080/AULm6Z",
	"original_url": "https://google.com",
	"created_at": "2026-09-27T11:49:42+03:00"
}
```

**Переход по ссылке:** открыть `http://localhost:8080/AULm6Z` в браузере — редирект на оригинальный URL.

**Статистика:**
```bash
curl http://localhost:8080/shorten/AULm6Z/stats
```

## Что было изучено на этом проекте

- Базовая работа с `net/http`: обработчики, `ResponseWriter`/`Request`, статус-коды
- Ручной разбор метода и параметров пути при роутинге через `http.HandleFunc`
- Проектирование через интерфейсы: `Storage` отделяет бизнес-логику от конкретного хранилища — переключение между map и PostgreSQL не потребовало менять `handlers` вообще
- Dependency injection через конструктор (`NewHandler`, `NewMemoryStorage`, `NewPostgresStorage`)
- Работа с JSON в теле HTTP-запроса и ответа
- Работа с `database/sql`: `QueryRow`/`Exec`, параметризованные запросы (`$1, $2`) для защиты от SQL-инъекций, `RowsAffected()` для различения «не найдено» от ошибки, `errors.Is(err, sql.ErrNoRows)`
- Переменные окружения и `.env` для хранения креды вне кода/репозитория
- Разделение на `internal/` (код приложения) и публичные пакеты (`models`, `dto`)
- Рекурсивная проверка коллизии при генерации уникального ID

## Известные ограничения / что улучшить дальше

- `sync.RWMutex` в `MemoryStorage` пока не задействован — планируется при изучении конкурентности
- Нет автоматических тестов
- Нет пула соединений с тонкой настройкой (`SetMaxOpenConns` и т.д.) — используются значения по умолчанию

## Автор

flecc1 — учебный проект в рамках изучения Go