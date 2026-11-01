# Календарь звонков

[![hexlet-check](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/one-thousand-twoth/ai-for-developers-project-386/actions)

Сервис для записи на звонки: владелец календаря публикует доступные интервалы
для встреч, а другой человек выбирает свободное время и бронирует встречу.

Учебный проект Хекслета: https://ru.hexlet.io/programs/ai-for-developers
Как это должно работать: https://files.hexlet.app/a/2ipc5m

## Текущее состояние

Сейчас реализована демонстрация создания бронирования и его получения по идентификатору
через API и веб-интерфейс. Данные хранятся в памяти и теряются при перезапуске сервера;
полные сценарии целевого продукта пока не реализованы.

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
Переменные окружения для демонстрации не нужны.

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
