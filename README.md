# Календарь звонков


[![hexlet-check](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions)

Разработайте совместно с ИИ сервис для бронирования календаря

Учебный проект Хекслета: https://ru.hexlet.io/programs/ai-for-developers
Как это должно работать: https://files.hexlet.app/a/2ipc5m

## Стек

- Backend: Go
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
mise run run
# В другом терминале:
mise run web:dev
```

Backend слушает http://localhost:8080, Vite — обычно http://localhost:5173.
Vite проксирует `/api` и `/health` на backend. Форма принимает местное время
устройства и отправляет RFC 3339 с часовым поясом через сгенерированный SDK.
В production frontend и API должны быть доступны с одного origin (или через reverse proxy).
Переменные окружения для демонстрации не нужны.

## Контракт и генерация

Единственный источник контракта — `typespec/main.tsp`.
Цепочка по образцу [Sphericon](https://github.com/getsphericon/sphericon/blob/main/docs/agents/codegen.md):

```text
typespec/ → OpenAPI 3.0 JSON в openapi/ → ogen → gen/api/
                                     → Hey API → web/src/generated/api/
```

TypeSpec и Hey API закреплены в корневом `package.json`, ogen — Go tool dependency
в `go.mod` (`go tool ogen`). Корневой pnpm workspace включает `web` и использует один lockfile.
SDK содержит TypeScript-типы и fetch-клиент; даты в JSON и SDK остаются строками.

После изменения контракта:

```bash
mise run generate
mise run generate:check
```

`generate` последовательно компилирует TypeSpec и запускает оба генератора.
`generate:typespec`, `generate:backend`, `generate:frontend` доступны отдельно;
последние две команды используют существующий `openapi/openapi.json`.
Коммитьте все артефакты вместе с контрактом. **Не редактируйте вручную**
`openapi/`, `gen/api/`, `web/src/generated/api/`: они перезаписываются генераторами
и помечены в `.gitattributes`. Проверка генерации обнаруживает изменения,
удаления и новые файлы только в этих каталогах; перед проверкой новые артефакты
должны быть добавлены в Git. Посторонние изменения рабочего дерева не мешают проверке.

## HTTP-примеры

```bash
curl http://localhost:8080/health
# {"status":"ok"}
curl -i http://localhost:8080/api/bookings \
  -H 'Content-Type: application/json' \
  -d '{"title":"Созвон","startsAt":"2026-11-01T15:00:00+08:00","endsAt":"2026-11-01T15:30:00+08:00"}'
# 201; скопируйте id из ответа:
curl http://localhost:8080/api/bookings/REPLACE_WITH_ID
```

Название обязательно и должно содержать непробельный символ. Время — RFC 3339
с часовым поясом, конец строго позже начала. Сервер возвращает даты в UTC.
Неверный запрос получает `400`, неизвестный ID — `404`; ошибки имеют вид
`{"code":"bad_request","message":"..."}` или `{"code":"not_found","message":"..."}`.

Хранилище находится в памяти процесса, защищено от конкурентного доступа,
ID — UUID. После перезапуска все бронирования исчезают. База данных, авторизация,
календари, участники и проверка пересечений времени в демонстрацию не входят.

## Проверки

```bash
mise run generate:check
mise run test
mise exec -- go test -race ./...
mise run lint
mise run build
mise run web:lint
mise run web:build
```

`web:install` устанавливает весь workspace с frozen lockfile.
`web:preview` показывает сборку; dev proxy работает только в `web:dev`.
Отдельный workflow `codegen.yml` проверяет генерацию, Go и frontend.
HTTP-тесты проверяют health, создание/чтение, ошибки и конкурентные запросы.

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
