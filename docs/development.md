# Разработка

Все команды выполняются из корня репозитория. [Установка](../README.md#установка-и-запуск).

## Переводы строк

Все текстовые файлы проекта используют LF. Правило `* text=auto eol=lf`
в `.gitattributes` сохраняет LF при checkout даже с `core.autocrlf=true`,
а `.editorconfig` задаёт LF в редакторах. Бинарные файлы не преобразуются.
Проверка: `git ls-files --eol` — текстовые файлы должны иметь `i/lf` и `w/lf`
(либо `none`, если переводов строк нет).

## Запуск и конфигурация

```bash
CALENDAR_TIMEZONE=Europe/Moscow mise run run
mise run web:dev # в другом терминале
```

`CALENDAR_TIMEZONE` обязателен: именованный IANA-пояс (`Europe/Moscow`, `Etc/UTC`, `UTC`).
`Local` и числовые смещения не подходят. База по умолчанию — `./data/calendar.sqlite`
(`DATABASE_PATH`), адрес — `:8080` (`LISTEN_ADDR`). Пример — `.env.example`;
сервер не загружает `.env` автоматически. Ошибки конфигурации, открытия базы и миграций
прерывают запуск до открытия HTTP-порта. Реальный адрес печатается после `net.Listen`;
`LISTEN_ADDR=127.0.0.1:0` позволяет ОС атомарно выбрать свободный порт.

`GET /` и `GET /health` возвращают `200` с JSON `{"status":"ok"}`.
Корневой маршрут нужен для проверки готовности контейнера в Hexlet CI.

Vite проксирует `/api` и `/health` на backend. В production нужен один origin
или reverse proxy. Старый UI демонстрирует только бронирования в памяти:
они теряются при перезапуске и не являются целевой записью гостей.
Типы встреч сохраняются в SQLite. Новый UI в этот этап не входит.

## Контракт и генерация

```text
typespec/ → openapi/ → ogen → gen/api/
                    → Hey API → packages/sdk/src/generated/
internal/storage/schema/ → Ent → internal/storage/ent/
```

TypeSpec — единственный источник HTTP-контракта. Версии TypeSpec/Hey API
закреплены в `package.json`, ogen/Ent/Atlas — в `go.mod` и `go.sum`.
`@calendar/sdk` — отдельный workspace-пакет без React, используемый UI и HTTP-тестами.
`mise run sdk:build` компилирует SDK в JavaScript и декларации типов в игнорируемом
`packages/sdk/dist/`. API-тесты и UI импортируют публичные exports `@calendar/sdk`;
тесты не зависят от внутреннего пути генерации. `test:api`, `web:dev` и `web:build`
автоматически собирают SDK перед запуском.

```bash
mise run generate
git add openapi gen/api packages/sdk/src/generated internal/storage/ent
mise run generate:check
```

Артефакты генерируются, а не редактируются вручную. Проверка drift обнаруживает
изменённые, удалённые и неотслеживаемые файлы, а также нестабильный повтор генерации.
Операции `generate:typespec`, `generate:backend`, `generate:frontend`, `generate:storage`
доступны отдельно. Генерация не создаёт миграций.

## SQLite и миграции

Ent-схемы отделены от API-моделей. Драйвер modernc.org/sqlite работает без CGO.
Соединения включают внешние ключи, ожидание блокировки и `_txlock=immediate`;
создание выполняется через один транзакционный клиент Ent. Уникальный индекс
защищает нормализованный Unicode-ключ и при нескольких процессах.

После изменения схемы и `mise run generate:storage`:

```bash
NAME=create_meeting_types mise run storage:diff
```

Команда Ent/Atlas воспроизводит предыдущие миграции в отдельной временной SQLite
и записывает SQL и `atlas.sum` в `internal/storage/migrations/`. Проверьте и закоммитьте
их вместе со схемой. При запуске сервер проверяет контрольные суммы и применяет
только ещё не применённые версии в IMMEDIATE-транзакции, сохраняя версии и хеши
в `schema_migrations`. Расхождение истории и ошибка SQL прерывают запуск с откатом.
`Schema.Create` при запуске не используется. Для резервной копии остановите сервер
и скопируйте файл базы.

### Ошибка `migration checksum: checksum mismatch`

Проверка относится к SQL-файлам миграций и `atlas.sum`, а не к содержимому SQLite.
Git с `core.autocrlf=true` может заменить LF на CRLF при checkout и тем самым
изменить контрольные суммы SQL. Для миграций в `.gitattributes` закреплены LF.
Проверить рабочую копию можно командой:

```bash
git ls-files --eol internal/storage/migrations
```

Ожидается `w/lf`. Если показано `w/crlf`, сохраните локальные изменения и получите
новую копию репозитория с актуальными атрибутами: обновление только `.gitattributes`
не переписывает уже выгруженные файлы. Пересоберите сервер или Docker-образ.
Удалять базу или пересоздавать `atlas.sum` ради обхода этой ошибки не нужно.
Если файлы уже в LF, проверьте, не были ли изменены сами миграции.

## HTTP-примеры

```bash
curl http://localhost:8080/health
# {"status":"ok"}
curl http://localhost:8080/api/meeting-types
# []
curl -i http://localhost:8080/api/owner/meeting-types \
  -H 'Content-Type: application/json' \
  -d '{"name":"Знакомство","durationMinutes":30,"description":"Первый звонок"}'
# 201; {"id":"...","name":"Знакомство","durationMinutes":30,"description":"Первый звонок"}
```

Название: 1–200 Unicode-символов и хотя бы один непробельный, пробелы по краям
удаляются при сохранении. Для уникальности используется Unicode case folding:
`Знакомство` и ` знакомство ` совпадают, внутренние пробелы значимы.
Длительность — целое 1–480 минут. Описание необязательно, до 2000 символов,
сохраняется без изменений; отсутствие и пустая строка различаются.
Ограничения длины применяются к переданной строке до удаления пробелов.

Ошибки — JSON `code/message`: `400 bad_request`,
`409 duplicate_meeting_type`, `503 temporarily_unavailable` при блокировке базы.
Неожиданные ошибки — `500 internal_error`, подробности только в логах.

## Проверки

```bash
mise run generate:check
mise run test
mise exec -- go test -race ./...
mise run lint
mise run build
mise run test:api
mise run web:lint
mise run web:build
```

`test:api` использует `tsc`, `node:test`, `node:assert/strict`; собирает настоящий
сервер с `CGO_ENABLED=0`, запускает его на loopback-порту, выбранном ОС, и ждёт health
с конечным таймаутом. Каждый тест получает отдельную базу в игнорируемом `.test/api/`.
Harness умеет остановить процесс и запустить его с той же базой; процессы и данные
очищаются и при ошибке, логи выводятся при сбое. После `mise run sdk:build` компиляция тестов — отдельно:
`mise exec -- pnpm exec tsc -p tests/api/tsconfig.json`.
Один файл после сборки: `mise exec -- node --test .test/api/build/tests/api/meeting-types.test.js`.

CI `codegen.yml` проверяет генерацию, Go/race, линт, сборку, SDK HTTP-сценарии и frontend.
Защищённый `hexlet-check.yml` не изменяется.

[Вся документация](README.md)
