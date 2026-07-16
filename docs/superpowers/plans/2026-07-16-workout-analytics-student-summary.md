# Workout Analytics + Student Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans or superpowers:subagent-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enriched workout detail (plan/fact + activity) for athlete and coach; week/month/year summary stats on student profile (coach) and progress (athlete).

**Architecture:** Extend existing DTOs/usecases/screens. No new routes. Period stats computed in usecases from `FindByUser` activities + plan workouts. Workout `Get` loads linked activity via `activityRepo.GetByID`.

**Tech Stack:** Go usecases + mockery/testify; React/TS screens; existing ActivityDetail UI pieces.

**Spec:** `docs/superpowers/specs/2026-07-16-workout-analytics-student-summary-design.md`

## Global Constraints

- Errors wrapped with `fmt.Errorf("path: %w", err)`
- Tests: no `MatchedBy` — pass concrete objects
- `planCompletion` with 0 planned workouts → `"—"`
- No new dedicated screens
- `make mock` after interface changes; `go test` green

## File map

### Create
- `apps/api/internal/domain/dto/period_stats.go` — PeriodStats + Summary
- `apps/api/internal/pkg/periodstats/periodstats.go` — build week/month/year from activities+workouts
- `apps/api/internal/usecase/student_usecase/get_test.go`
- `apps/api/internal/usecase/workout_usecase/get_test.go` (or extend create_update_test)
- `apps/web/src/components/stats/PeriodSummary.tsx`
- `apps/web/src/components/workout/WorkoutFactPanel.tsx`

### Modify
- `apps/api/internal/domain/dto/workout_view.go` — PlannedKm, ActualKm, ActualPace, Fact
- `apps/api/internal/domain/dto/student_detail.go` — Summary, PlanDay.WorkoutID
- `apps/api/internal/domain/dto/activity_view.go` — ProgressResponse.Summary
- `apps/api/internal/adapter/postgres/workout_repo/repo.go` — FindByUser (non-template)
- `apps/api/internal/usecase/workout_usecase/*` — Get enrich + activityRepo.GetByID
- `apps/api/internal/usecase/student_usecase/*` — summary, actualKm, workoutId
- `apps/api/internal/usecase/activity_usecase/progress.go` — summary
- Web: `api.ts`, `WorkoutDetailScreen`, `StudentProfileScreen`, `ProgressScreen`

---

### Task 1: DTOs + periodstats helper + FindByUser

- [ ] Add `PeriodStatsView` / `PeriodSummaryView` in `period_stats.go`
- [ ] Extend WorkoutView, PlanDayView, StudentDetailView, ProgressResponse
- [ ] Add `workout_repo.FindByUser` (is_club_template=false)
- [ ] Implement `periodstats.Build(now, activities, planWorkouts) PeriodSummaryView`
- [ ] Unit-test periodstats (table test) OR cover via usecase tests
- [ ] Commit

### Task 2: Workout Get + fact

- [ ] Extend activityRepo interface with `GetByID`
- [ ] Get: verify coach membership; attach Fact from activity; set ActualKm/Pace
- [ ] Update returns enriched view too
- [ ] Tests + `make mock` + `go test ./internal/usecase/workout_usecase/...`
- [ ] Commit

### Task 3: Student detail summary

- [ ] planDays include workoutId; actualKm from linked activity dist
- [ ] summary week/month/year via periodstats
- [ ] Extend mocks/interfaces; tests
- [ ] Commit

### Task 4: Progress summary

- [ ] Progress adds Summary via periodstats (need workoutRepo FindByUser or FindByUserWeek×N + own)
- [ ] Tests
- [ ] Commit

### Task 5: Web UI

- [ ] Types in api.ts
- [ ] PeriodSummary component
- [ ] WorkoutDetail: plan/fact + inline activity panel (streams by fact.id)
- [ ] StudentProfile: period switcher + open workout
- [ ] Progress: period summary
- [ ] `npx tsc --noEmit`
- [ ] Commit

## Success check

- Coach: student → week/month/year → tap plan day → plan/fact + activity
- Athlete: completed workout detail + Progress summary
- Incomplete workout: empty fact, no crash
