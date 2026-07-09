# DONE_WITH_CONCERNS

Task 8 complete with simplifications:

- `comp` uses all-time `SumDistByUser` vs `plan_weeks[week_index=0].TargetKm()` (first number in `plan_label`); no calendar week window / date-based current week yet.
- Repo lives at `adapter/postgres/plan_week_repo` (project convention), not `internal/repository/`.
- Plan clamp uses min/max `week_index` from `FindByClub`; empty plan weeks → empty range/plan strings, workouts still returned.
