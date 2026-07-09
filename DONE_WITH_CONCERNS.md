# DONE_WITH_CONCERNS

Task 11 complete.

- BottomNav hybrid C: 5 tabs + SVG icons, floating `.nav` capsule; schedule not in nav (open via HomeAthlete `group-preview`).
- App uses `user.role === 'coach'`; no role switcher in-app.
- Hardcoded demo numbers stripped from HomeAthlete / HomeCoach / ProgressScreen / ProfileScreen.
- Playwright e2e updated for RoleEntry path + schedule via home card.
- Vitest + build green.

## Removed hardcoded fallbacks

| Location | Removed | Now |
|---|---|---|
| ProgressScreen | `?? 612`, `?? 78`, `?? 3` | API or `—` |
| HomeAthlete | `24.6` / `из 28 км` | sum of activities + `plan.weekPlan` |
| HomeCoach | `3` (сегодня), `86%` | `announces.length`, `analytics.attendance` or `—` |
| ProfileAthlete | `НП`, `97.6 кг`, тренер «Зина» | initials from `user.name`; fake rows removed |
| ProfileCoach | `Зина · ЛЭМЗ` row | removed |
| App avatar | hardcoded `НП`/`ТР` | initials from `user.name` |

## Concerns

- HomeAthlete «эта неделя» sums all loaded activities (API has no week-window filter) — may overcount vs true week km.
- HomeCoach middle card is announce count (not “today sessions”) — no today-session API.
- e2e not run here (needs compose/API); unit tests + build only.
- Avatar click → profile is a small UX extra; schedule still only from home card.
