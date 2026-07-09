# DONE_WITH_CONCERNS

Task 10 complete.

- Path A: RoleEntry → AuthGate(role) → CreateClubGate when `user.needsClub`.
- AuthGate: role from prop, onBack clears role, empty email/password (no seed prefill).
- CreateClubGate uses store `createClub` (name + accent), then refresh via store.
- Vitest + build green; BottomNav left for Task 11.

Concerns:
- `authRole` is local App state (not store) — fine for Path A; resets on full remount.
- Login still accepts any role path; backend role comes from JWT after auth (role prop only affects register).
- Athlete with `needsClub` would also hit CreateClubGate if API ever sets that flag for athletes.
