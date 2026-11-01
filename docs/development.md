# Разработка

[Установка и быстрый запуск](../README.md#установка-и-запуск) описаны в корневом README.
Все команды ниже выполняются из корня репозитория.

## Окружение и frontend

Vite проксирует `/api` и `/health` на backend. Форма принимает местное время
устройства и отправляет RFC 3339 с часовым поясом через сгенерированный SDK.
В production frontend и API должны быть доступны с одного origin (или через reverse proxy).
Переменные окружения для демонстрации не нужны.

`mise run web:install` устанавливает весь workspace с frozen lockfile.
`mise run web:preview` показывает сборку; dev proxy работает только в `mise run web:dev`.

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

Примеры относятся к текущей демонстрации, а не к [целевому продукту](product.md).

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
ID — UUID. После перезапуска все бронирования исчезают. В текущей демонстрации
ещё нет базы данных, типов встреч, календаря доступных слотов, участников
и проверки пересечений времени. Авторизация не входит в границы проекта.

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

Отдельный workflow `codegen.yml` проверяет генерацию, Go и frontend.
HTTP-тесты проверяют health, создание/чтение, ошибки и конкурентные запросы.

[Вся документация](README.md)
