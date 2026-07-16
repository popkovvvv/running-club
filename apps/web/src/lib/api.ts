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

export type Student = {
  id: string
  init: string
  name: string
  sub: string
  km: string
  comp: number
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
  removeStudent: (id: string) => request<{ status: string }>(`/club/students/${id}`, { method: 'DELETE' }),
  announces: () => request<Announce[]>('/announces'),
  publishAnnounce: (body: { place: string; day: string; time: string; group: string; note: string }) =>
    request<Announce>('/announces', { method: 'POST', body: JSON.stringify(body) }),
  signup: (id: string) => request<Announce>(`/announces/${id}/signup`, { method: 'POST' }),
  unsignup: (id: string) => request<Announce>(`/announces/${id}/signup`, { method: 'DELETE' }),
  calendar: () => request<{ cells: Array<{ key: string; n: number; blank: boolean; has: boolean; isToday: boolean; bg: string; fg: string; dot: string }> }>('/schedule/calendar'),
  plan: (week = 1) => request<{ weekIndex: number; weekRange: string; weekPlan: string; days: Array<{ id: string; dayLabel: string; tag: string; title: string; distKm: number; duration: string }>; mine: Array<{ id: string; dayLabel: string; tag: string; title: string; duration: string }> }>(`/plan?week=${week}`),
  createWorkout: (body: Record<string, unknown>) =>
    request('/workouts', { method: 'POST', body: JSON.stringify(body) }),
  activities: () => request<Activity[]>('/activities'),
  stravaIntegration: () => request<IntegrationStatus>('/integrations/strava'),
  connectStrava: () => request<{ url: string }>('/integrations/strava/connect'),
  disconnectStrava: () => request<{ status: string }>('/integrations/strava/disconnect', { method: 'POST' }),
  progress: () => request<{ months: Array<{ m: string; km: number; tr: number; pace: string; diff: string }>; yearKm: number; yearTr: number; yearStarts: number }>('/progress'),
  analytics: () => request<{ clubKm: number; attendance: number; students: Student[] }>('/analytics'),
  races: () => request<Array<{ days: string; name: string; date: string; dist: string; goal: string }>>('/races'),
  prs: () => request<Array<{ d: string; t: string; date: string }>>('/prs'),
}
