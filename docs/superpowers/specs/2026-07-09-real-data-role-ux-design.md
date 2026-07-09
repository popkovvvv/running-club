# Real data + role UX

**Date:** 2026-07-09  
**Status:** approved for planning  
**Approach:** domain/data first, then UI

## Goals

1. Remove hardcoded/mock presentation data — UI reads from Postgres via API.
2. Bottom nav matches design hybrid: 5 tabs + icons, keep floating capsule; schedule not in nav.
3. Tabs and screens split strictly by account role (`user.role`).
4. Auth entry: choose athlete or coach; coach gets a dedicated create-club step after register.

## Decisions locked

| Topic | Choice |
|-------|--------|
| Bottom nav | Hybrid C: 5 tabs, SVG icons, floating capsule; schedule via home card |
| Auth entry | Path A: role first, then login/register |
| Coach club | Separate step after coach registration (`POST /clubs`) |
| Role switcher in header | Remove completely |
| Seed accounts | Local only (`SEED=1`); no UI prefill |
| Implementation order | Domain/API first, then UI |

## Out of scope

- Payments / subscriptions
- Changing role of an existing user
- Full calendar CRUD editor
- Dual-theme (STRIDE); PULSE only

---

## 1. Domain and data

### Remove hardcodes

| Current | Target |
|---------|--------|
| `progress` year totals `612/78/3` | `yearKm` = sum(`month_stats.km`); `yearTr` = sum(`month_stats.tr`); `yearStarts` = count of user's `races` |
| `analytics` `186.1`, `86%`, student `"24.6"` | `clubKm` = sum of club athletes' activity `dist_km`; `attendance` = round(100 * signed_up_slots / max(going_capacity,1)) across club announces (or 0 if no announces); per-student `km` = sum of that athlete's activities |
| `list_students` stub `sub/km/comp` | `km` = athlete activity sum (formatted); `sub` = next upcoming announce label for club or «Нет записи»; `comp` = min(100, round(100 * athlete_week_km / plan_km)) for current `plan_weeks` row, else 0 |
| `calendar` fixed July 2026 + session days | Current month grid (`time.Now()`); `has` on days with club `announces` (match by day-of-month from announce `day_label` or add `starts_on date` if needed); accent from club |
| `weekMeta` in Go code | Persist as `plan_weeks` |
| Register coach auto-creates club | Register creates user only; club via `POST /clubs` |

### Schema

1. New table `plan_weeks`:
   - `id` UUID PK (app-generated)
   - `club_id` UUID NOT NULL REFERENCES clubs(id)
   - `week_index` INT NOT NULL
   - `range_label` TEXT NOT NULL
   - `plan_label` TEXT NOT NULL
   - UNIQUE (`club_id`, `week_index`)

2. `announces.starts_on DATE NULL` — used by calendar to mark days; seed fills real dates; UI publish should send/store date (day_label can remain display text).

Migrations: create only via `make migration` / `oh-my-pg-tool goose create`, native goose SQL (`-- +goose Up/Down` + `StatementBegin/End`).

### API changes

- `POST /api/v1/clubs` — body `{ name, accentHex }`; coach-only; fails if club already exists for coach; returns `ClubView` with invite code.
- Auth/`me` response: `needsClub: true` when `role=coach` and no club for coach.
- `Register` for coach: **no** club insert.
- Calendar handler: pass club accent (not hardcoded `#ff5c22`).
- Progress / analytics / students / calendar / plan week meta: repository-backed only.

### Seed

- Runs only when `SEED=1` (local compose/dev).
- Keeps demo users/club/announces/activities for local testing.
- Creates `plan_weeks` rows for seed club.
- Must not prefill AuthGate fields.

---

## 2. Auth and navigation UX

### Auth flow

1. Gate: **Я спортсмен** / **Я тренер**.
2. Then login or register (role fixed from step 1; no role toggle on the form).
3. Athlete after register → app (join club by code remains existing flow).
4. Coach after register → if `needsClub` → **Создать клуб** (name, accent) → `POST /clubs` → app.
5. Login of coach without club → same create-club gate.

### In-app role

- All screen/nav branching uses `user.role` only.
- Remove `demoRole` / `setDemoRole` and seed re-login switcher.

### Bottom nav (hybrid C)

Tabs:

| key | Athlete | Coach |
|-----|---------|-------|
| home | Главная | Ученики |
| plan | План | Конструктор |
| prog | Прогресс | Аналитика |
| races | Старты | Старты |
| profile | Профиль | Клуб |

- Icons + labels like legacy.
- Keep current floating capsule styles.
- Schedule: open from athlete home card (not a nav item).

### Frontend cleanup

- Remove fallbacks like `?? 612`, hardcoded home stats (`24.6`, `86%`, etc.) — show API values or empty/`—`.
- Update Vitest/Playwright for new auth paths and 5-tab nav.

---

## 3. Delivery order

### Phase 1 — domain/API

1. Migration `plan_weeks` (+ any repo methods needed for aggregates).
2. Repos/usecases: progress, analytics, students, calendar, plan weeks from DB.
3. Register without auto-club; `POST /clubs`; `needsClub` on auth/me.
4. Seed update under `SEED=1`.
5. Unit tests (mockery) + API e2e for create club / calendar / progress.

### Phase 2 — UI

1. Role entry → login/register → CreateClubGate.
2. Remove role switcher and credential prefill.
3. Bottom nav hybrid C; schedule from home.
4. Screens keyed by `user.role`; no mock number fallbacks.

### Phase 3 — verify

- `make migrate-up`, local `SEED=1`
- Vitest + Playwright (athlete vs coach tabs, coach create club)
- Compose smoke: login via web `:8088` → API without 502

---

## 4. Testing notes

- Usecase tests: mockery only; no hand-written repo mocks; no `MatchedBy`; no `tt := tt`.
- Errors wrapped with `fmt.Errorf("path: %w", err)` as project convention.
- Domain entities via constructors in `internal/domain/model`; DTO views via constructors in `internal/domain/dto`.

## 5. Success criteria

- No hardcoded year/analytics/calendar/student metrics in usecases.
- Coach cannot enter main app without a club; create-club step works end-to-end.
- Athlete and coach see different nav labels/screens from the same build, based on JWT user role.
- Bottom nav has 5 capsule tabs with icons; schedule reachable from home.
- Seed data available only when explicitly enabled for local test.
