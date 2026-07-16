# Running Club

Мобильный веб клуба бега: React + Vite (тема **PULSE 1a**) и Go/chi API. Деплой на Vercel, локально — Docker Compose или dev-режим.

## Быстрый старт (Podman)

```bash
podman machine start          # если машина ещё не запущена
cp .env.example .env
podman compose up --build
```

- Web: http://localhost:8088
- API: http://localhost:8080/healthz
- Seed: `nikita@pulse.run` / `password` (спортсмен), `coach@pulse.run` / `password` (тренер), код клуба `PULSE-7K42`

Остановка: `podman compose down`

Если порт `5432` занят локальным Postgres (`postgresql@16`), либо останови его, либо в `docker-compose.yml` смени проброс на `"5433:5432"`.

## Dev без полного compose

Только БД в Podman, api/web локально:

```bash
podman compose up postgres -d
cp .env.example .env
cd apps/api && HTTP_ADDR=:18080 SEED=1 go run ./cmd/api
cd apps/web && npm install && npm run dev
```

Web: http://localhost:5173 (proxy `/api` → `:18080`)

## Make

| Команда | Описание |
|---------|----------|
| `make up` / `make down` | `podman compose` |
| `make migration` | новая SQL-миграция через `oh-my-pg-tool goose create` |
| `make migrate-up` / `make migrate-down` | накатить / откатить через `oh-my-pg-tool local` |
| `make migrate-status` | статус goose-миграций |
| `make test-unit` | Go unit (`-tags=unit`) |
| `make test-e2e` | Go API e2e (нужен Postgres) |
| `make test-web` | Vitest |
| `make test-web-e2e` | Playwright |
| `make seed` | демо-данные |

Миграции — нативный SQL goose в `apps/api/scripts/migrations/` (шаблон `-- +goose Up/Down` + `StatementBegin/End`). Создавать только через `make migration` (`oh-my-pg-tool`). Нужен `oh-my-pg-tool` в `PATH`. DSN по умолчанию: `postgres://pulse:pulse@localhost:5432/running_club?sslmode=disable` (переопределяется `MIGRATE_DSN`).

Для e2e API создайте БД:

```bash
podman compose exec postgres psql -U pulse -c 'CREATE DATABASE running_club_test;'
TEST_DATABASE_URL=postgres://pulse:pulse@localhost:5432/running_club_test?sslmode=disable make test-e2e
```

## Структура

```
apps/web   — React UI
apps/api   — Go API (domain → usecase → adapter → http)
legacy/    — исходный .dc.html прототип
```

## Vercel

Один проект, два сервиса (раздельно собираются):

| Service | Root | Runtime |
|---|---|---|
| `web` | `apps/web` | Vite (static) |
| `api` | `apps/api` | Go (`cmd/api`, слушает `PORT`) |

1. Импортируйте репозиторий (Root Directory = корень репо)
2. Env: `DATABASE_URL` (Neon), `JWT_SECRET`, при необходимости Strava
3. Роутинг: `/api/*` → api, остальное → web (см. `vercel.json`)

Нужен доступ к **Vercel Services** в аккаунте/плане. Локально: `vercel dev` из корня.

Если Services недоступны — напишите, сделаем fallback на классический `api/index.go` + отдельный build web.

## Что закрыто из прототипа

- Без оплат/абонементов
- Членство: join по коду, leave, invite code, remove student
- `scheduleCta`: «Записаться» / «Вы записаны»
- Blank-ячейки календаря без точки
- Один вариант UI — PULSE 1a + палитра клуба
