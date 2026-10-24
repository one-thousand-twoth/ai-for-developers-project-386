# Calendar

Go-модуль: `github.com/one-thousand-twoth/calendar`. Версия Go указана в `go.mod`.

## Структура

- `cmd/server` — точка входа HTTP-сервера.
- `modules` — модули приложения.
- `config` — конфигурация приложения.
- `internal` — внутренние детали реализации.
- `docs` — общая документация.
- `web` — пустая папка для будущего фронтенда Vite + React.

## Запуск

```sh
go run ./cmd/server
```

Сервер слушает порт `8080`. Проверка состояния:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Сборка

```sh
go build -o /tmp/calendar-server ./cmd/server
```

GitHub Actions выполняет smoke-проверку сборкой `./cmd/server`.

Git не хранит пустые папки: после клонирования при необходимости создайте `web`
командой `mkdir -p web`.
