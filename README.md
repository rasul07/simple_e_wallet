# E-Wallet

## О проекте

E-Wallet - это проект электронного кошелька, разработанный с использованием Go. Проект предоставляет API для управления электронными кошельками пользователей, включая операции с балансом и аутентификацию.

## Структура проекта

- `api/`: Документация API (Swagger)
- `cmd/`: Точка входа в приложение
- `internal/`: Внутренняя логика приложения
  - `config/`: Конфигурация приложения
  - `handlers/`: Обработчики HTTP-запросов
  - `models/`: Модели данных
  - `service/`: Бизнес-логика
  - `storage/`: Взаимодействие с базой данных
- `migrations/`: SQL-миграции для базы данных
- `Dockerfile`: Инструкции для сборки Docker-образа
- `docker-compose.yml`: Конфигурация для запуска приложения и зависимостей
- `go.mod` и `go.sum`: Зависимости проекта
- `Makefile`: Команды для упрощения разработки и запуска

## Требования

- Docker
- Docker Compose

## Запуск проекта

1. Клонируйте репозиторий:
   ```
   git clone <URL репозитория>
   cd alif_ewallet
   ```

2. Запустите проект с помощью Docker Compose:
   ```
   docker-compose up --build
   ```

   Это действие выполнит следующее:
   - Соберет Docker-образ приложения
   - Запустит контейнер с PostgreSQL
   - Применит миграции к базе данных
   - Запустит приложение на порту 8080

3. После успешного запуска, API будет доступно по адресу: `http://localhost:8080`

## Разработка

- Для добавления новых миграций, создайте SQL-файл в директории `migrations/`
- Используйте `Makefile` для выполнения часто используемых команд (если они определены)

## Тестирование

Для запуска тестов используйте команду:
```
go test ./...
```

## Документация API

Swagger-документация доступна в директории `api/docs/`. После запуска приложения, она может быть доступна через эндпоинт `/swagger` (если настроено).

## Переменные окружения

Приложение использует следующие переменные окружения:
- `DB_PORT`: Порт базы данных (по умолчанию 5432)
- `DB_USERNAME`: Имя пользователя базы данных
- `DB_PASSWORD`: Пароль пользователя базы данных
- `DB_NAME`: Имя базы данных
- `DB_SSL_MODE`: Режим SSL для подключения к базе данных

Эти переменные можно настроить в файле `docker-compose.yml` или передать напрямую при запуске приложения.
