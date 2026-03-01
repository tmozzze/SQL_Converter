# SQL Converter API

## Описание проекта
Данный сервис предоставляет API для загрузки файлов форматов `.csv` и `.xlsx`. Программа автоматически анализирует содержимое файла, определяет типы данных для каждой колонки (Integer, Float, Boolean, String) и создает таблицу в PostgreSQL.

## Запуск проекта

### Предварительные требования
- Установленный Docker и Docker Compose.
- Наличие файла `.env` в корне проекта (создайте его на основе примера `.env.example`).

### Быстрый старт
1. **Соберите и запустите контейнеры. Для этого в терминали bash:**
   ```sh
   make up
   ```
   *Эта команда поднимет базу данных PostgreSQL и сам API сервер.*

2. **Откройте Swagger-документацию:**
   ```
   http://localhost:8080/swagger/index.html
   ```

3. **Тестирование:**
   В корне проекта находятся тестовые файлы: `test.csv`, `test2.csv`, `test3.xlsx`. Вы можете загрузить их через Swagger UI или cURL.

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go            # Точка входа в приложение
├── config/
│   └── local.yaml             # Конфигурация для локальной разработки
├── docs/                      # Сгенерированная Swagger документация
├── internal/
│   ├── config/                # Загрузка конфига
│   ├── domain/                # Доменные сущности
│   │   ├── models/            # Доменные модели
│   │   ├── errors.go          # Ошибки
│   │   ├── repository.go      # Интерфейсы слоя данных
│   │   └── service.go         # Интерфейсы бизнес-логики
│   ├── http/
│   │   └── handler/           # Хендлеры
│   ├── repository/
│   │   └── postgres/          # Слой репозитория
│   └── service/               # Реализация бизнес-логики
│       ├── analyzer.go        # Алгоритм определения типов данных
│       ├── parser.go          # Парсинг CSV и XLSX
│       └── processor.go       # Управление процессом загрузки
├── pkg/
│   └── database/              # Подключение к БД
├── Dockerfile                 
├── docker-compose.yml         
├── Makefile                   # Команды для управления проектом
└── go.mod                     # Зависимости
```

## Технологический стек
- **Language:** Go 1.26
- **Database:** PostgreSQL 15
- **HTTP Framework:** `net/http`
- **Logging:** `log/slog`
- **Config:** `cleanenv`, `godotenv`
- **Documentation:** Swagger (`swaggo`)


## Команды управления (Makefile)

### Разработка
- `make test` — запуск тестов
- `make lint` — проверка кода линтером `golangci-lint`
- `make swagger-gen` — перегенерация документации Swagger

### Docker
- `make up` — сборка и запуск проекта в Docker.
- `make down` — остановка и удаление контейнеров.
- `make down-and-clean` — полная очистка (удаление контейнеров и данных БД).
