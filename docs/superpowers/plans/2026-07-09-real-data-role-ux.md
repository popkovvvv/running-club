# Real Data + Role UX Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Persist presentation data in Postgres, add coach create-club flow, and ship role-based auth + hybrid bottom nav without demo role switching.

**Architecture:** Phase 1 extends domain/repos/usecases (migrations via `oh-my-pg-tool`, constructors in `model`/`dto`, mockery tests). Phase 2 rewires React auth/nav to `user.role` and removes hardcoded UI fallbacks. Seed stays behind `SEED=1` only.

**Tech Stack:** Go/chi/pgx/goose, React/Vite/TS, mockery + testify, Vitest/Playwright, Podman compose.

**Spec:** `docs/superpowers/specs/2026-07-09-real-data-role-ux-design.md`

---

## File map

### Create
- `apps/api/scripts/migrations/<ts>_plan_weeks_and_announce_starts_on.sql`
- `apps/api/internal/domain/model/plan_week.go`
- `apps/api/internal/domain/dto/create_club.go`
- `apps/api/internal/adapter/postgres/plan_week_repo/repo.go`
- `apps/api/internal/usecase/club_usecase/create.go`
- `apps/api/internal/usecase/club_usecase/create_test.go`
- `apps/web/src/components/RoleEntry.tsx`
- `apps/web/src/components/CreateClubGate.tsx`
- `apps/web/src/components/BottomNav.tsx`

### Modify (API)
- `apps/api/internal/domain/model/announce.go` — `StartsOn *time.Time`
- `apps/api/internal/domain/dto/user_view.go` — `NeedsClub`
- `apps/api/internal/domain/dto/announce_view.go` / create announce request — `startsOn`
- `apps/api/internal/usecase/auth_usecase/register.go` — no auto club
- `apps/api/internal/usecase/auth_usecase/helpers.go` — `needsClub` for coach
- `apps/api/internal/usecase/club_usecase/usecase.go` — create deps
- `apps/api/internal/usecase/club_usecase/list_students.go` — real metrics
- `apps/api/internal/usecase/activity_usecase/progress.go` — aggregates
- `apps/api/internal/usecase/analytics_usecase/club_analytics.go` — aggregates
- `apps/api/internal/usecase/schedule_usecase/calendar.go` + `publish.go`
- `apps/api/internal/usecase/workout_usecase/plan.go` + `usecase.go` — plan weeks from DB
- `apps/api/internal/adapter/postgres/*` — new queries
- `apps/api/internal/app/http/club/handler.go` + `router/router.go`
- `apps/api/internal/pkg/seed/seed.go`
- `apps/api/.mockery.yaml` — regenerate mocks as interfaces change

### Modify (Web)
- `apps/web/src/components/AuthGate.tsx`
- `apps/web/src/lib/api.ts`, `store.tsx`
- `apps/web/src/App.tsx`
- screens using `demoRole` / hardcoded numbers
- Vitest + Playwright specs

---

### Task 1: Migration `plan_weeks` + `announces.starts_on`

**Files:**
- Create: `apps/api/scripts/migrations/<timestamp>_plan_weeks_and_announce_starts_on.sql`

- [ ] **Step 1: Create migration via oh-my-pg-tool**

```bash
cd /Users/nikpopkov/GolandProjects/running-club
# non-interactive equivalent of make migration:
oh-my-pg-tool goose create plan_weeks_and_announce_starts_on -d apps/api/scripts/migrations
```

- [ ] **Step 2: Fill SQL**

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE announces ADD COLUMN starts_on DATE;

CREATE TABLE plan_weeks (
    id UUID PRIMARY KEY,
    club_id UUID NOT NULL REFERENCES clubs(id),
    week_index INT NOT NULL,
    range_label TEXT NOT NULL,
    plan_label TEXT NOT NULL,
    UNIQUE (club_id, week_index)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_weeks;
ALTER TABLE announces DROP COLUMN IF EXISTS starts_on;
-- +goose StatementEnd
```

- [ ] **Step 3: Apply locally**

```bash
make migrate-up
make migrate-status
```

Expected: new migration applied; status shows both `00001_init.sql` and the new file.

- [ ] **Step 4: Commit**

```bash
git add apps/api/scripts/migrations
git commit -m "Add plan_weeks table and announces.starts_on"
```

---

### Task 2: Domain model + DTO for plan weeks and needsClub / create club

**Files:**
- Create: `apps/api/internal/domain/model/plan_week.go`
- Modify: `apps/api/internal/domain/model/announce.go`
- Modify: `apps/api/internal/domain/dto/user_view.go`
- Create: `apps/api/internal/domain/dto/create_club.go`
- Modify: `apps/api/internal/domain/dto/announce_view.go` (CreateAnnounceRequest)

- [ ] **Step 1: Add `PlanWeek` constructor**

```go
package model

import "github.com/google/uuid"

type PlanWeek struct {
	ID         uuid.UUID
	ClubID     uuid.UUID
	WeekIndex  int
	RangeLabel string
	PlanLabel  string
}

func NewPlanWeek(clubID uuid.UUID, weekIndex int, rangeLabel, planLabel string) *PlanWeek {
	return &PlanWeek{
		ID: uuid.New(), ClubID: clubID, WeekIndex: weekIndex,
		RangeLabel: rangeLabel, PlanLabel: planLabel,
	}
}
```

- [ ] **Step 2: Add `StartsOn` to `Announce` + `NewAnnounce` param `startsOn *time.Time`**

Update `NewAnnounce` signature to accept `startsOn *time.Time` and set `a.StartsOn = startsOn`. Update all call sites (publish, seed, tests/fixtures).

- [ ] **Step 3: Extend UserView**

```go
type UserView struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	InClub    bool       `json:"inClub"`
	NeedsClub bool       `json:"needsClub"`
	ClubID    *uuid.UUID `json:"clubId,omitempty"`
}

func (v *UserView) MarkNeedsClub() *UserView {
	v.NeedsClub = true
	v.InClub = false
	return v
}
```

- [ ] **Step 4: CreateClub DTO**

```go
package dto

type CreateClubRequest struct {
	Name      string `json:"name"`
	AccentHex string `json:"accentHex"`
}
```

Add `StartsOn string \`json:"startsOn,omitempty"\`` (`YYYY-MM-DD`) to `CreateAnnounceRequest`.

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/domain
git commit -m "Add plan week model and needsClub/create-club DTOs"
```

---

### Task 3: Auth — register without auto-club + needsClub

**Files:**
- Modify: `apps/api/internal/usecase/auth_usecase/register.go`
- Modify: `apps/api/internal/usecase/auth_usecase/helpers.go`
- Modify: `apps/api/internal/usecase/auth_usecase/register_test.go`
- Modify: `apps/api/internal/usecase/auth_usecase/usecase.go` (clubCreator may become unused — remove dep if unused)
- Regenerate mocks if interfaces change: `cd apps/api && mockery`

- [ ] **Step 1: Write failing register coach test (no club create)**

In `register_test.go`, case `ok_coach`:

```go
{
  name: "ok_coach_no_auto_club",
  req:  dto.RegisterRequest{Name: "Тренер", Email: "c@b.c", Password: "secret1", Role: "coach"},
  before: func(m usecaseMocks) {
    m.userRepo.EXPECT().GetByEmail(mock.Anything, "c@b.c").Return(nil, model.ErrNotFound).Once()
    m.userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
    // clubRepo.Create must NOT be expected
    m.clubRepo.EXPECT().GetByCoachID(mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()
  },
},
```

Assert `res.User.NeedsClub == true`, `res.User.InClub == false`.

- [ ] **Step 2: Run test — expect fail (still auto-creates club / MarkInClub)**

```bash
cd apps/api && go test -tags=unit ./internal/usecase/auth_usecase -run TestRegister -count=1
```

- [ ] **Step 3: Implement**

`register.go`: delete coach club creation block.

`helpers.go` `toView`:

```go
if user.Role == model.RoleCoach {
  club, err := u.clubRepo.GetByCoachID(ctx, user.ID)
  if err == nil {
    return view.WithClub(club.ID), nil
  }
  if !errors.Is(err, model.ErrNotFound) {
    return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
  }
  return view.MarkNeedsClub(), nil
}
```

Change `clubCreator` interface to include `GetByCoachID` **or** inject `clubRepo` with Create+GetByCoachID. Prefer renaming to `clubRepo` with both methods. Update `NewUseCase`, service provider wiring, `.mockery.yaml`, run `mockery`.

- [ ] **Step 4: Tests green**

```bash
cd apps/api && go test -tags=unit ./internal/usecase/auth_usecase -count=1
```

- [ ] **Step 5: Commit**

```bash
git add apps/api/internal/usecase/auth_usecase apps/api/.mockery.yaml apps/api/internal/usecase/auth_usecase/mocks apps/api/cmd
git commit -m "Stop auto-creating club on coach register; expose needsClub"
```

---

### Task 4: Create club usecase + HTTP

**Files:**
- Modify: `apps/api/internal/usecase/club_usecase/usecase.go` — add `Create` to clubRepo if missing
- Create: `apps/api/internal/usecase/club_usecase/create.go`
- Create: `apps/api/internal/usecase/club_usecase/create_test.go`
- Modify: `apps/api/internal/app/http/club/handler.go`
- Modify: `apps/api/internal/app/http/router/router.go`
- Modify: postgres `club_repo` if `Create` missing
- `make mock` after interface changes

- [ ] **Step 1: Failing test `TestCreate`**

```go
func TestCreate(t *testing.T) {
  t.Parallel()
  coachID := uuid.New()
  tests := []struct {
    name string
    before func(m usecaseMocks)
    wantErr error
  }{
    {
      name: "ok",
      before: func(m usecaseMocks) {
        m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(nil, model.ErrNotFound).Once()
        m.clubRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Club")).Return(nil).Once()
        m.clubRepo.EXPECT().CountActiveStudents(mock.Anything, mock.Anything).Return(0, nil).Once()
      },
    },
    {
      name: "already_has_club",
      before: func(m usecaseMocks) {
        m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(model.ClubFixture(uuid.New(), "X", "X", "#ff5c22", coachID), nil).Once()
      },
      wantErr: model.ErrConflict,
    },
  }
  // ...
}
```

- [ ] **Step 2: Run — fail (method missing)**

```bash
cd apps/api && go test -tags=unit ./internal/usecase/club_usecase -run TestCreate -count=1
```

- [ ] **Step 3: Implement `Create`**

```go
func (u *UseCase) Create(ctx context.Context, coachID uuid.UUID, req dto.CreateClubRequest) (*dto.ClubView, error) {
  name := strings.TrimSpace(req.Name)
  accent := strings.TrimSpace(req.AccentHex)
  if name == "" {
    return nil, fmt.Errorf("%w: empty name", model.ErrInvalidInviteCode)
  }
  if !strings.HasPrefix(accent, "#") || len(accent) != 7 {
    accent = "#ff5c22"
  }
  if _, err := u.clubRepo.GetByCoachID(ctx, coachID); err == nil {
    return nil, model.ErrConflict
  } else if !errors.Is(err, model.ErrNotFound) {
    return nil, fmt.Errorf("clubRepo.GetByCoachID: %w", err)
  }
  code := "PULSE-" + strings.ToUpper(coachID.String()[:4])
  club := model.NewClub(name, code, accent, coachID)
  if err := u.clubRepo.Create(ctx, club); err != nil {
    return nil, fmt.Errorf("clubRepo.Create: %w", err)
  }
  n, _ := u.clubRepo.CountActiveStudents(ctx, club.ID)
  return dto.NewClubView(club.ID, club.Name, club.AccentHex, n).WithInviteCode(club.InviteCode), nil
}
```

Ensure `clubRepo` interface has `Create`. Wire handler:

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
  if middleware.Role(r.Context()) != string(model.RoleCoach) {
    response.Error(w, model.ErrForbidden)
    return
  }
  var req dto.CreateClubRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    response.Error(w, err)
    return
  }
  res, err := h.uc.Create(r.Context(), middleware.UserID(r.Context()), req)
  // ...
}
```

Router: `priv.Post("/clubs", h.Club.Create)`.

- [ ] **Step 4: Tests + build**

```bash
cd apps/api && mockery && go test -tags=unit ./internal/usecase/club_usecase -count=1 && go build ./cmd/api
```

- [ ] **Step 5: Commit**

```bash
git add apps/api
git commit -m "Add POST /clubs create-club flow for coaches"
```

---

### Task 5: Progress aggregates from DB

**Files:**
- Modify: `apps/api/internal/usecase/activity_usecase/progress.go`
- Modify: `apps/api/internal/usecase/activity_usecase/progress_test.go`
- Possibly extend `activityRepo` with nothing if month stats + races already available

- [ ] **Step 1: Failing test**

Mock `FindMonthStats` returning two months (km 100+50, tr 10+5) and `FindRaces` returning 2 races → expect `yearKm=150`, `yearTr=15`, `yearStarts=2`.

- [ ] **Step 2: Run — fail on hardcoded 612**

- [ ] **Step 3: Implement**

```go
var yearKm float64
var yearTr int
for _, m := range months {
  yearKm += m.Km
  yearTr += m.Tr
  views = append(views, dto.NewMonthStatView(...))
}
races, err := u.activityRepo.FindRaces(ctx, userID)
// wrap error
yearStarts := len(races)
return dto.NewProgressResponse(yearKm, yearTr, yearStarts, views), nil
```

- [ ] **Step 4: Tests green + commit**

```bash
cd apps/api && go test -tags=unit ./internal/usecase/activity_usecase -count=1
git add apps/api/internal/usecase/activity_usecase
git commit -m "Compute progress year totals from month_stats and races"
```

---

### Task 6: Students + analytics from DB

**Files:**
- Modify: `apps/api/internal/usecase/club_usecase/list_students.go` + tests
- Modify: `apps/api/internal/usecase/analytics_usecase/*`
- Extend repos: activity sums by user/club, announce signup counts — add methods on existing repos or inject activity/announce repos into club/analytics usecases

**Repo methods to add (postgres + interfaces):**
- `activityRepo.SumDistByUser(ctx, userID) (float64, error)`
- `activityRepo.SumDistByClubAthletes(ctx, clubID) (float64, error)`
- `announceRepo.AttendanceStats(ctx, clubID) (signedUp, capacity int, err error)` — capacity = sum going_count or signup seats; use `sum(going_count)` and `count(signups)` as defined in spec
- `announceRepo.NextForAthlete(ctx, clubID, athleteID) (label string, err error)`
- `planWeekRepo.GetByClubWeek` / `FindByClub`
- `workoutRepo` or activities for week km — use sum of athlete activities in current week window **or** sum `workouts` own+plan dist for current week_index; prefer activities `dist_km` for athlete week km

- [ ] **Step 1: Extend interfaces + mockery + failing tests for ListStudents / ClubAnalytics**

- [ ] **Step 2: Implement repo SQL + usecases per spec formulas**

- [ ] **Step 3:**

```bash
cd apps/api && mockery && go test -tags=unit ./internal/usecase/club_usecase ./internal/usecase/analytics_usecase -count=1
git commit -m "Derive student and club analytics metrics from DB"
```

---

### Task 7: Calendar from announces + club accent

**Files:**
- Modify: `apps/api/internal/usecase/schedule_usecase/calendar.go`
- Modify: `apps/api/internal/usecase/schedule_usecase/publish.go`
- Modify: `apps/api/internal/app/http/schedule/handler.go`
- Modify: announce postgres repo DAO for `starts_on`
- Tests: `calendar_test.go`

- [ ] **Step 1: Change Calendar signature**

```go
func (u *UseCase) Calendar(ctx context.Context, userID uuid.UUID, role string) (*dto.CalendarResponse, error)
```

Resolve club via existing `clubIDFor`; load announces; build set of day-of-month from `StartsOn` in current month; use `time.Now()` for month/today; load club accent via clubRepo (add `GetByID` to schedule deps if needed) for colors.

- [ ] **Step 2: Handler**

```go
res, err := h.uc.Calendar(r.Context(), middleware.UserID(r.Context()), middleware.Role(r.Context()))
```

- [ ] **Step 3: Publish stores StartsOn**

Parse `req.StartsOn` (`2006-01-02`); if empty, leave nil.

- [ ] **Step 4: Tests + commit**

```bash
cd apps/api && go test -tags=unit ./internal/usecase/schedule_usecase -count=1
git commit -m "Build schedule calendar from announce dates and club accent"
```

---

### Task 8: Plan weeks from DB + seed update

**Files:**
- Create: `apps/api/internal/adapter/postgres/plan_week_repo/repo.go`
- Modify: `apps/api/internal/usecase/workout_usecase/usecase.go` — remove `weekMeta` const; inject `planWeekRepo`
- Modify: `apps/api/internal/usecase/workout_usecase/plan.go`
- Modify: `apps/api/internal/pkg/seed/seed.go`
- Wire in service provider / main DI

- [ ] **Step 1: plan_week_repo Create/FindByClub/GetByClubAndIndex**

- [ ] **Step 2: Plan usecase loads weeks for athlete's club (via membership) or empty defaults**

If no plan weeks: return empty range/plan strings and still return workouts.

- [ ] **Step 3: Seed inserts plan_weeks for seed club + set announce `starts_on`**

- [ ] **Step 4:**

```bash
cd apps/api && mockery && go test -tags=unit ./internal/usecase/workout_usecase -count=1
SEED=1 go run ./cmd/api # smoke migrate+seed if DB up
git commit -m "Load plan week metadata from plan_weeks; update seed"
```

---

### Task 9: Web API client + store — needsClub, createClub, drop demoRole

**Files:**
- Modify: `apps/web/src/lib/api.ts`
- Modify: `apps/web/src/lib/store.tsx`

- [ ] **Step 1: Types**

```ts
export type User = {
  id: string
  name: string
  email: string
  role: 'athlete' | 'coach'
  inClub: boolean
  needsClub?: boolean
  clubId?: string
}

// api.createClub = (body: { name: string; accentHex: string }) =>
//   request<Club>('/clubs', { method: 'POST', body: JSON.stringify(body) })
```

- [ ] **Step 2: Remove `demoRole` / `setDemoRole` from store; use `user.role`**

- [ ] **Step 3: Commit**

```bash
git add apps/web/src/lib
git commit -m "Add createClub client and remove demo role state"
```

---

### Task 10: Auth UI — role entry + create club gate

**Files:**
- Create: `apps/web/src/components/RoleEntry.tsx`
- Create: `apps/web/src/components/CreateClubGate.tsx`
- Modify: `apps/web/src/components/AuthGate.tsx`
- Modify: `apps/web/src/App.tsx`
- Update: `AuthGate.test.tsx` + new component tests

- [ ] **Step 1: RoleEntry** — two buttons set pending role in local state lifted to App/store (`authRole: 'athlete'|'coach'|null`)

- [ ] **Step 2: AuthGate** — receive `role` prop; remove role toggle; empty email/password defaults; register uses prop role

- [ ] **Step 3: CreateClubGate** — name + accent inputs; submit `api.createClub`; then `refresh()`

- [ ] **Step 4: App gating**

```tsx
if (!user) {
  if (!authRole) return <RoleEntry onPick={setAuthRole} />
  return <AuthGate role={authRole} onBack={() => setAuthRole(null)} />
}
if (user.needsClub) return <CreateClubGate />
```

- [ ] **Step 5: Vitest**

```bash
cd apps/web && npm test -- --run
git commit -m "Add role entry and create-club auth gates"
```

---

### Task 11: Bottom nav hybrid C + screens by user.role

**Files:**
- Create: `apps/web/src/components/BottomNav.tsx`
- Modify: `apps/web/src/App.tsx`
- Modify: `apps/web/src/styles/app.css` (nav button icon layout if needed)
- Modify screens: replace `demoRole` with `user.role`
- Remove hardcoded numbers in HomeAthlete/HomeCoach/ProgressScreen/ProfileScreen
- Update Playwright `e2e/app.spec.ts`

- [ ] **Step 1: BottomNav** — 5 tabs with SVG icons (copy paths from `legacy/PhoneApp.dc.html`), labels by `user.role`, keep `.nav` capsule class; no schedule tab

- [ ] **Step 2: App** — remove role switcher UI; `isCoach = user.role === 'coach'`; schedule screen still reachable via `setScreen('schedule')` from HomeAthlete card

- [ ] **Step 3: Strip UI fallbacks** (`?? 612`, `24.6`, `86%`, etc.)

- [ ] **Step 4:**

```bash
cd apps/web && npm test -- --run
# if API up:
cd apps/web && npm run test:e2e
git commit -m "Ship hybrid bottom nav and role-based screens without mocks"
```

---

### Task 12: End-to-end verification

- [ ] **Step 1: API unit**

```bash
cd apps/api && go test -tags=unit ./...
```

- [ ] **Step 2: Migrate + compose smoke**

```bash
cd /Users/nikpopkov/GolandProjects/running-club
make migrate-up
podman compose up -d --build
curl -sS http://127.0.0.1:8080/healthz
curl -sS -X POST http://127.0.0.1:8088/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"nikita@pulse.run","password":"password"}'
```

Expected: health `ok`; login via web proxy `200` (not 502).

- [ ] **Step 3: Manual checklist**
  - Register coach → create club → coach tabs
  - Register athlete → join code → athlete tabs (5), schedule from home
  - Progress/analytics numbers match DB seed aggregates

- [ ] **Step 4: Final commit if fixes needed**

```bash
git commit -m "Fix verification issues for real-data role UX"
```

---

## Spec coverage checklist

| Spec item | Task |
|-----------|------|
| plan_weeks + starts_on migration | 1 |
| needsClub / no auto club on register | 3 |
| POST /clubs | 4 |
| progress aggregates | 5 |
| students + analytics aggregates | 6 |
| calendar from DB + accent | 7 |
| week meta from DB + seed | 8 |
| seed local only / no prefill | 8, 10 |
| role entry auth A | 10 |
| create club step A | 4, 10 |
| remove role switcher | 9, 11 |
| bottom nav hybrid C | 11 |
| screens by user.role | 11 |
| remove UI hardcode fallbacks | 11 |
| verify compose/login | 12 |
