# Workout analytics + student summary for coach

**Date:** 2026-07-16  
**Status:** approved for planning  
**Approach:** expand existing WorkoutDetail / StudentProfile / Progress screens; enrich current API endpoints (no new standalone screens)

## Goals

1. Show plan vs fact and full activity analytics for a specific workout (athlete and coach).
2. Show summary statistics for a student to the coach: week / month / year (km, workouts, plan completion %, avg pace).
3. Give the athlete the same period summary on Progress.

## Decisions locked

| Topic | Choice |
|-------|--------|
| Workout analytics depth | C: plan/fact summary + full activity breakdown on one screen |
| Student summary horizon | A: week + month + year |
| Who sees workout analytics | B: coach and athlete |
| Architecture | Expand existing screens + enrich existing endpoints |
| planCompletion with no planned workouts | Show `—` (not `0%`) |
| New dedicated analytics screens | Out of scope |

## Out of scope

- Separate «Analytics workout» / «Student stats» screens
- Client-only stitching without API changes
- Athlete-only later / coach-only later rollouts
- Changing Strava sync itself
- Club-wide period charts beyond current analytics list

---

## 1. Data and API

### `GET /workouts/{id}` (enriched)

Keep current plan fields. Add optional fact block when a linked activity exists (`completedActivityId` / activity `linkedWorkoutId`):

| Field group | Contents |
|-------------|----------|
| Plan (existing) | title, type, day, distKm, pace, status, segments, description |
| Comparison | `plannedKm`, `actualKm`, planned pace vs actual pace (actual from activity) |
| Fact (optional) | same core metrics as ActivityDetail: dist, time, pace, HR avg/max, elev, kudos, polyline, splits |
| Streams | unchanged: client still loads `GET /activities/{id}/streams` when fact is present |

If no linked activity: return plan only; UI shows empty fact state.

Access: workout owner (athlete) or coach of the athlete’s club; otherwise project-standard 403/404.

### `GET /club/students/{id}` (coach summary)

Add:

```text
summary: {
  week:  PeriodStats
  month: PeriodStats
  year:  PeriodStats
}
```

`PeriodStats`: `km`, `workouts` (completed count), `planCompletion` (string `%` or `—`), `avgPace` (string or `—`).

Also fix/enrich `planDays[]`:

- add `workoutId`
- `actualKm` from linked activity distance when completed (not planned dist)

Week selector behavior stays; summary periods are calendar week / calendar month / calendar year relative to «now» (not the week picker). Document in plan if week picker should later drive week summary — **v1: summary is always current periods**.

### `GET /progress` (athlete)

Add the same `summary: { week, month, year }` shape so Progress UI mirrors coach student header.

Existing year totals + monthly table remain.

### Period definitions

| Period | Window |
|--------|--------|
| week | current calendar week (Mon–Sun), consistent with existing week km helpers |
| month | current calendar month |
| year | current calendar year |

`planCompletion` = round(100 * completed_planned_workouts / planned_workouts) for workouts scheduled in that window; if `planned_workouts == 0` → `—`.

`avgPace` = average pace over activities in the window that have pace; if none → `—`.

---

## 2. UI and navigation

### WorkoutDetail (athlete + coach)

1. Existing plan header + segments.
2. Block «План / факт»: planned vs actual km and pace; status.
3. If activity linked: reuse ActivityDetail sections inline (map, metrics, splits, streams).
4. If not linked: «Ещё нет данных»; athlete keeps «Отметить выполненной».
5. Separate ActivityDetail overlay remains available from feeds; not required for workout path.

### StudentProfile (coach)

1. Top: period switcher Неделя | Месяц | Год with `PeriodStats`.
2. Existing week plan list: rows open WorkoutDetail via `workoutId`.
3. Recent activities unchanged (open activity detail).

### Progress (athlete)

1. Same period summary switcher above existing year/month content.
2. Own plan workouts continue to open enriched WorkoutDetail.

---

## 3. Access, empty states, tests

### Access

- Enriched workout: owner or club coach.
- Student summary: coach of that club only.
- Progress summary: authenticated athlete self only.

### Empty states

- No linked activity → plan UI only; no crash.
- No streams/map → hide those sections.
- `planCompletion` / `avgPace` → `—` when undefined.

### Tests

- Usecase: workout get/update complete returns fact when activity exists.
- Usecase: student detail includes `workoutId`, real `actualKm`, and three period summaries.
- Usecase: progress includes summary periods.
- Mocks via mockery; tests build concrete args (no `MatchedBy`).
- After API changes: `make mock` and `go test` for touched packages.
- Web: API types + screens consume new fields; e2e optional.

---

## 4. Implementation order

1. Backend DTOs + student/workout/progress usecases + tests.
2. Web API types.
3. WorkoutDetail enrichment (reuse activity UI pieces).
4. StudentProfile summary + clickable plan days.
5. Progress athlete summary.

## Success criteria

- Coach opens student → sees week/month/year stats → taps a workout → sees plan/fact + activity analytics.
- Athlete opens completed workout from plan → same workout analytics; Progress shows period summary.
- Incomplete workout still usable; no broken empty states.
