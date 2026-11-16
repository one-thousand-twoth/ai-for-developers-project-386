# Календарь звонков

[![hexlet-check](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions)

Сервис для записи на звонки: владелец календаря публикует доступные интервалы
для встреч, а другой человек выбирает свободное время и бронирует встречу.

Учебный проект Хекслета: https://ru.hexlet.io/programs/ai-for-developers
Как это должно работать: https://files.hexlet.app/a/2ipc5m

## Текущее состояние

Реализованы создание и чтение типов встреч через API и независимый TypeScript SDK.
Типы сохраняются в SQLite после перезапуска. Старый веб-интерфейс остаётся демонстрацией
бронирований в памяти; слоты, запись гостей и управление встречами ещё не реализованы.

## Стек

- Backend: Go, Ent, SQLite (modernc, без CGO)
- Frontend: React, TypeScript, Vite
- Зависимости фронтенда: pnpm
- Управление инструментами и задачами: mise

## Установка и запуск

Версии Go, Node.js, pnpm и golangci-lint закреплены в `mise.toml`.
Установите [mise](https://mise.jdx.dev/), затем из корня проекта:

```bash
mise trust
mise install
mise run install
CALENDAR_TIMEZONE=Europe/Moscow mise run run
# В другом терминале:
mise run web:dev
```

Backend слушает http://localhost:8080, Vite — обычно http://localhost:5173.
`CALENDAR_TIMEZONE` обязателен: именованный IANA-пояс, например `Europe/Moscow` или `Etc/UTC`.
`DATABASE_PATH` по умолчанию `./data/calendar.sqlite`, `LISTEN_ADDR` — `:8080`.
Каталог базы создаётся автоматически, версионированные миграции применяются при запуске.
Ошибка конфигурации, открытия базы или миграции останавливает сервер.

Проверка API без React: `mise run test:api` (tsc, node:test, настоящий сервер и отдельная SQLite).

### Backend в Docker

```bash
docker build -t calendar-backend .
docker run --rm -p 8080:8080 \
  -v calendar-data:/data \
  calendar-backend
```

API доступен на http://localhost:8080. Образ содержит только бекенд и запускается
без root. SQLite хранится в томе `calendar-data`, который сохраняется после удаления
контейнера. В образе задан `CALENDAR_TIMEZONE=Etc/UTC`; другой пояс можно указать
через `-e CALENDAR_TIMEZONE=Europe/Moscow` при запуске.

## Документация

- [Целевой продукт](docs/product.md) — сценарии и бизнес-правила.
- [Разработка](docs/development.md) — генерация API, HTTP-примеры и проверки.
- [Глоссарий](GLOSSARY.md) — термины проекта.
- [Вся документация](docs/README.md), включая выпуск релизов.

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
