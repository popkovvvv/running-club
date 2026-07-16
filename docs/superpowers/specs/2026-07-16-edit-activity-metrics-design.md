# Edit activity metrics — Design Spec

**Date:** 2026-07-16  
**Status:** Approved for planning

## Problem

After marking a workout completed, the linked manual activity often has distance from the plan but empty time, pace, and HR. Athletes and coaches also need to correct imported (e.g. Strava) metrics locally. There is no API or UI to edit activity fields.

## Goals

- Allow editing activity fact data from both the activity detail screen and the workout detail screen (when a linked `fact` exists).
- Support all activity sources (`manual`, `strava`, etc.).
- Athlete can edit own activities; coach can edit activities of active club students.
- Editable fields: title, date, distance (km), duration, pace, heart rate, elevation gain.

## Non-goals

- Pushing edits back to Strava or other providers.
- Editing streams, polyline/map, kudos, or comments.
- Changing workout plan segments via this flow.
- Deleting activities.

## Approach

**`PATCH /activities/{id}`** plus a shared edit form used on activity and workout screens.

### Why not extend workout PATCH only

Strava (and other) activities may exist without a workout link. Editing through activities covers every entry point.

## API

### `PATCH /activities/{id}`

**Auth:** required.

**Authorization:**
- actor is activity owner, OR
- actor is coach of a club where the activity owner has an active membership.

**Request body** (all fields optional; only sent fields are updated):

| Field | Type | Notes |
|-------|------|--------|
| `title` | string | |
| `when` | string | `YYYY-MM-DD`; sets `startedAt` and regenerates `whenLabel` |
| `distKm` | number | also updates `distanceMeters` |
| `duration` | string | display time, e.g. `46:30` |
| `pace` | string | e.g. `5:15` |
| `hr` | number | average HR; updates `hr` / `averageHeartrate` as used by views |
| `elevationGain` | number | meters |

**Response:** full `ActivityDetailView` (same shape as `GET /activities/{id}`).

**Errors:** `404` if not found; `403` if not owner/coach; `400` on invalid date/values.

### Side effects

- `updatedAt` bumped.
- Streams / polyline / kudos untouched.
- Linked workout (if any) is not mutated; `GET /workouts/{id}` already loads `fact` from the activity, so fact metrics refresh automatically.

## UI

### Shared `ActivityEditForm`

Bottom sheet / modal with fields: title, date, km, duration, pace, HR, elevation (m). Actions: Save / Cancel. On success: refresh local activity state and call `reloadActivities`.

### Entry points

1. **`ActivityDetailScreen`** — «Редактировать» under the header for users with access (owner or coach viewing student activity).
2. **`WorkoutFactPanel`** (workout detail with `fact`) — same button above metrics.

Coach opening a student workout/activity uses the same form and the same PATCH.

### Unchanged

- Metric grid remains read-only outside edit mode.
- Map and stream tabs stay view-only.

## Access rules (frontend)

- Athlete: show edit on own activities / own workout facts.
- Coach: show edit when viewing a student’s activity or workout fact (API enforces access).

## Testing

- Unit: activity usecase Update — owner ok; coach of member ok; outsider 403; partial patch; invalid date.
- Frontend: form submits PATCH and refreshes displayed metrics (component or screen-level as fits existing patterns).

## Out of scope follow-ups

- Auto-compute pace from dist + duration.
- Sync edits to Strava.
