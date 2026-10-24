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

## Установка

<!-- Опишите установку: клонирование, зависимости, переменные окружения -->

```bash
git clone https://github.com/one-thousand-twoth/ai-for-developers-project-386.git
cd ai-for-developers-project-386
```

## Фронтенд

Приложение находится в `web/`. Версии Node.js и pnpm зафиксированы в `mise.toml`.

Установите [mise](https://mise.jdx.dev/), затем выполните из корня проекта:

```bash
mise trust
mise install
mise run web:install
mise run web:dev
```

Vite выведет адрес приложения (обычно http://localhost:5173).

```bash
mise run web:build    # Проверка TypeScript и сборка в web/dist
mise run web:lint     # Проверка линтером
mise run web:preview  # Просмотр готовой сборки
```

Если mise активирован в оболочке, команды pnpm можно запускать напрямую из `web/`:

```bash
cd web
pnpm install --frozen-lockfile
pnpm dev
```

## Использование

<!-- Добавьте примеры запуска и запись asciinema — именно это смотрит работодатель -->

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
