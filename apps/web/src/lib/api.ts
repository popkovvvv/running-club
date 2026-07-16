import type { WorkoutType } from './workoutTypes'

const BASE = import.meta.env.VITE_API_BASE || '/api/v1'

export type User = {
  id: string
  name: string
  email: string
  role: 'athlete' | 'coach'
  inClub: boolean
  needsClub?: boolean
  clubId?: string
}

export type Club = {
  id: string
  name: string
  inviteCode?: string
  accentHex: string
  students: number
}

export type Announce = {
  id: string
  place: string
  day: string
  time: string
  group: string
  note: string
  going: number
  signedUp: boolean
  scheduleCta: string
  startsOn?: string
  attendees?: Array<{ id: string; init: string; name: string }>
}

export type Activity = {
  id: string
  title: string
  when: string
  dist: string
  time: string
  pace: string
  hr: string
  kudos: number
  comN: number
  route: string
  sx: number
  sy: number
  ex: number
  ey: number
  source?: string
  sportType?: string
  elevation?: number
  visibility?: string
}

export type ActivityDetail = Activity & {
  startedAt?: string
  movingSeconds?: number
  elapsedSeconds?: number
  maxHeartrate?: number
  polyline?: string
  externalId?: string
  linkedWorkoutId?: string
}

export type ActivityStream = { type: string; data: string }

export type Segment = { id?: string; kind: string; title: string; distKm: number; pace: string }

export type Workout = {
  id: string
  kind: string
  workoutType: WorkoutType
  dayLabel: string
  tag: string
  title: string
  description?: string
  distKm: number
  hr: string
  weekIndex: number
  scheduledDate?: string
  status: string
  completedActivityId?: string
  assignedBy?: string
  isClubTemplate?: boolean
  segments?: Segment[]
  plannedKm?: number
  actualKm?: number
  actualPace?: string
  rpe?: number
  athleteReport?: string
  coachComment?: string
  fact?: ActivityDetail
}

export type PlanWeek = { weekIndex: number; rangeLabel: string; planLabel: string }

export type PeriodStats = {
  km: string
  workouts: number
  planCompletion: string
  avgPace: string
}

export type PeriodSummary = {
  week: PeriodStats
  month: PeriodStats
  year: PeriodStats
}

export type Student = {
  id: string
  init: string
  name: string
  sub: string
  km: string
  comp: number
}

export type StudentDetail = {
  student: Student
  weekKm: string
  weekPlanKm: string
  comp: number
  year?: number
  month?: number
  planDays: Array<{
    workoutId: string
    dayLabel: string
    title: string
    workoutType: string
    plannedKm: number
    actualKm: number
    status: string
    scheduledDate?: string
    activityId?: string
    activityWhen?: string
    activityPace?: string
    rpe?: number
    athleteReport?: string
    coachComment?: string
  }>
  recentActivities: Activity[]
  summary: PeriodSummary
}

export type IntegrationStatus = {
  provider: string
  status: string
  connected: boolean
  externalAthleteId?: string
  scopes?: string[]
  expiresAt?: string
  lastSyncedAt?: string
  lastWebhookAt?: string
  lastError?: string
}

type AuthResponse = { token: string; user: User }

function authHeader(): HeadersInit {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...authHeader(),
      ...(init?.headers || {}),
    },
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(body.error || res.statusText)
  }
  return res.json()
}

export const api = {
  register: (body: { name: string; email: string; password: string; role: string }) =>
    request<AuthResponse>('/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: { email: string; password: string }) =>
    request<AuthResponse>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  me: () => request<User>('/auth/me'),
  logout: () => request<{ status: string }>('/auth/logout', { method: 'POST' }),
  club: () => request<Club>('/club'),
  createClub: (body: { name: string; accentHex: string }) =>
    request<Club>('/clubs', { method: 'POST', body: JSON.stringify(body) }),
  joinClub: (code: string) => request<Club>('/club/join', { method: 'POST', body: JSON.stringify({ code }) }),
  leaveClub: () => request<{ status: string }>('/club/leave', { method: 'POST' }),
  palette: (accentHex: string) =>
    request<Club>('/club/palette', { method: 'PATCH', body: JSON.stringify({ accentHex }) }),
  students: () => request<Student[]>('/club/students'),
  student: (id: string, opts?: { week?: number; year?: number; month?: number }) => {
    const q = new URLSearchParams()
    if (opts?.year) q.set('year', String(opts.year))
    if (opts?.month) q.set('month', String(opts.month))
    if (opts?.week != null) q.set('week', String(opts.week))
    const qs = q.toString()
    return request<StudentDetail>(`/club/students/${id}${qs ? `?${qs}` : ''}`)
  },
  removeStudent: (id: string) => request<{ status: string }>(`/club/students/${id}`, { method: 'DELETE' }),
  announces: () => request<Announce[]>('/announces'),
  publishAnnounce: (body: { place: string; day: string; time: string; group: string; note: string; startsOn?: string }) =>
    request<Announce>('/announces', { method: 'POST', body: JSON.stringify(body) }),
  signup: (id: string) => request<Announce>(`/announces/${id}/signup`, { method: 'POST' }),
  unsignup: (id: string) => request<Announce>(`/announces/${id}/signup`, { method: 'DELETE' }),
  calendar: (year?: number, month?: number) => {
    const q = new URLSearchParams()
    if (year) q.set('year', String(year))
    if (month) q.set('month', String(month))
    const suffix = q.toString() ? `?${q}` : ''
    return request<{
      year: number
      month: number
      label: string
      cells: Array<{ key: string; n: number; blank: boolean; has: boolean; isToday: boolean; iso?: string; bg: string; fg: string; dot: string }>
    }>(`/schedule/calendar${suffix}`)
  },
  plan: (opts?: number | { week?: number; year?: number; month?: number }) => {
    const q = new URLSearchParams()
    if (typeof opts === 'number') {
      q.set('week', String(opts))
    } else if (opts) {
      if (opts.week != null) q.set('week', String(opts.week))
      if (opts.year) q.set('year', String(opts.year))
      if (opts.month) q.set('month', String(opts.month))
    }
    const qs = q.toString()
    return request<{ weekIndex: number; weekRange: string; weekPlan: string; weekKm: string; days: Workout[]; mine: Workout[] }>(`/plan${qs ? `?${qs}` : ''}`)
  },
  planWeeks: () => request<PlanWeek[]>('/plan/weeks'),
  upsertPlanWeek: (weekIndex: number, body: { rangeLabel: string; planLabel: string }) =>
    request<PlanWeek>(`/plan/weeks/${weekIndex}`, { method: 'PUT', body: JSON.stringify(body) }),
  planTemplate: (weekIndex: number) => request<{ weekIndex: number; workouts: Workout[] }>(`/plan/weeks/${weekIndex}/template`),
  savePlanTemplate: (weekIndex: number, workouts: Record<string, unknown>[]) =>
    request<{ weekIndex: number; workouts: Workout[] }>(`/plan/weeks/${weekIndex}/template`, { method: 'PUT', body: JSON.stringify({ workouts }) }),
  publishPlan: (weekIndex: number) => request<{ status: string }>(`/plan/weeks/${weekIndex}/publish`, { method: 'POST' }),
  createWorkout: (body: Record<string, unknown>) => request<Workout>('/workouts', { method: 'POST', body: JSON.stringify(body) }),
  getWorkout: (id: string) => request<Workout>(`/workouts/${id}`),
  updateWorkout: (id: string, body: Record<string, unknown>) =>
    request<Workout>(`/workouts/${id}`, { method: 'PATCH', body: JSON.stringify(body) }),
  deleteWorkout: (id: string) => request<{ status: string }>(`/workouts/${id}`, { method: 'DELETE' }),
  activities: () => request<Activity[]>('/activities'),
  activity: (id: string) => request<ActivityDetail>(`/activities/${id}`),
  updateActivity: (id: string, body: Record<string, unknown>) =>
    request<ActivityDetail>(`/activities/${id}`, { method: 'PATCH', body: JSON.stringify(body) }),
  activityStreams: (id: string) => request<ActivityStream[]>(`/activities/${id}/streams`),
  stravaIntegration: () => request<IntegrationStatus>('/integrations/strava'),
  connectStrava: () => request<{ url: string }>('/integrations/strava/connect'),
  disconnectStrava: () => request<{ status: string }>('/integrations/strava/disconnect', { method: 'POST' }),
  progress: () => request<{ months: Array<{ m: string; km: number; tr: number; pace: string; diff: string }>; yearKm: number; yearTr: number; yearStarts: number; summary: PeriodSummary }>('/progress'),
  analytics: () => request<{ clubKm: number; students: Student[] }>('/analytics'),
  races: () => request<Array<{ days: string; name: string; date: string; dist: string; goal: string }>>('/races'),
  prs: () => request<Array<{ d: string; t: string; date: string }>>('/prs'),
}
