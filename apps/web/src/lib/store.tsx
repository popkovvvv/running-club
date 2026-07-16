import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { api, type Activity, type Announce, type Club, type User } from './api'
import { pulseTheme, withAccent, type Theme } from './theme'

type Screen = 'home' | 'schedule' | 'plan' | 'prog' | 'profile'

export type Overlay =
  | { type: 'activity'; id: string }
  | { type: 'workout'; id: string }
  | { type: 'student'; id: string }
  | null

type AppState = {
  user: User | null
  token: string | null
  theme: Theme
  club: Club | null
  screen: Screen
  tabEpoch: number
  overlay: Overlay
  announces: Announce[]
  activities: Activity[]
  loading: boolean
  setScreen: (s: Screen) => void
  openOverlay: (o: Overlay) => void
  closeOverlay: () => void
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string, role: string) => Promise<void>
  logout: () => void
  refresh: () => Promise<void>
  createClub: (body: { name: string; accentHex: string }) => Promise<void>
  joinClub: (code: string) => Promise<void>
  leaveClub: () => Promise<void>
  setPalette: (hex: string) => Promise<void>
  signupAnnounce: (id: string) => Promise<void>
  unsignupAnnounce: (id: string) => Promise<void>
  removeStudent: (id: string) => Promise<void>
  publishAnnounce: (body: { place: string; day: string; time: string; group: string; note: string; startsOn?: string }) => Promise<void>
  reloadAnnounces: () => Promise<void>
  reloadActivities: () => Promise<void>
}

const Ctx = createContext<AppState | null>(null)

export function AppProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [club, setClub] = useState<Club | null>(null)
  const [screen, setScreenState] = useState<Screen>('home')
  const [tabEpoch, setTabEpoch] = useState(0)
  const [overlay, setOverlay] = useState<Overlay>(null)
  const [announces, setAnnounces] = useState<Announce[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)

  const setScreen = useCallback((s: Screen) => {
    setScreenState(s)
    setTabEpoch((n) => n + 1)
  }, [])

  const closeOverlay = useCallback(() => {
    setOverlay(null)
    setTabEpoch((n) => n + 1)
  }, [])

  const theme = useMemo(
    () => ({
      ...withAccent(pulseTheme, club?.accentHex || pulseTheme.accent),
      name: club?.name || pulseTheme.name,
    }),
    [club?.accentHex, club?.name],
  )

  const reloadActivities = useCallback(async () => {
    if (!localStorage.getItem('token')) {
      setActivities([])
      return
    }
    try {
      const list = await api.activities()
      setActivities(list)
    } catch {
      setActivities([])
    }
  }, [])

  const refresh = useCallback(async () => {
    if (!localStorage.getItem('token')) {
      setUser(null)
      setClub(null)
      setAnnounces([])
      setActivities([])
      setLoading(false)
      return
    }
    try {
      const me = await api.me()
      setUser(me)
      if (me.needsClub) {
        setClub(null)
        setAnnounces([])
        setActivities([])
      } else if (me.inClub || me.role === 'coach') {
        const c = await api.club()
        setClub(c)
        setAnnounces(await api.announces())
        if (me.role !== 'coach') {
          setActivities(await api.activities().catch(() => []))
        } else {
          setActivities([])
        }
      } else {
        setClub(null)
        setAnnounces([])
        setActivities(await api.activities().catch(() => []))
      }
    } catch {
      localStorage.removeItem('token')
      setToken(null)
      setUser(null)
      setActivities([])
    } finally {
      setLoading(false)
    }
  }, [])

  const reloadAnnounces = useCallback(async () => {
    if (!localStorage.getItem('token')) {
      setAnnounces([])
      return
    }
    try {
      setAnnounces(await api.announces())
    } catch {
      setAnnounces([])
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh, token])

  const applyAuth = async (res: { token: string; user: User }) => {
    localStorage.setItem('token', res.token)
    setToken(res.token)
    setUser(res.user)
    setScreen('home')
    await refresh()
  }

  const value: AppState = {
    user,
    token,
    theme,
    club,
    screen,
    tabEpoch,
    overlay,
    announces,
    activities,
    loading,
    setScreen,
    openOverlay: setOverlay,
    closeOverlay,
    login: async (email, password) => {
      const res = await api.login({ email, password })
      await applyAuth(res)
    },
    register: async (name, email, password, role) => {
      const res = await api.register({ name, email, password, role })
      await applyAuth(res)
    },
    logout: () => {
      localStorage.removeItem('token')
      setToken(null)
      setUser(null)
      setClub(null)
      setAnnounces([])
      setActivities([])
    },
    refresh,
    createClub: async (body) => {
      const c = await api.createClub(body)
      setClub(c)
      setUser((u) => (u ? { ...u, needsClub: false, inClub: true, clubId: c.id } : u))
      await refresh()
    },
    joinClub: async (code) => {
      const c = await api.joinClub(code)
      setClub(c)
      setAnnounces(await api.announces())
      setUser((u) => (u ? { ...u, inClub: true, clubId: c.id } : u))
      await reloadActivities()
    },
    leaveClub: async () => {
      await api.leaveClub()
      setClub(null)
      setAnnounces([])
      setUser((u) => (u ? { ...u, inClub: false, clubId: undefined } : u))
      setScreen('home')
    },
    setPalette: async (hex) => {
      const c = await api.palette(hex)
      setClub(c)
    },
    signupAnnounce: async (id) => {
      const a = await api.signup(id)
      setAnnounces((list) => list.map((x) => (x.id === id ? a : x)))
    },
    unsignupAnnounce: async (id) => {
      const a = await api.unsignup(id)
      setAnnounces((list) => list.map((x) => (x.id === id ? a : x)))
    },
    removeStudent: async (id) => {
      await api.removeStudent(id)
    },
    publishAnnounce: async (body) => {
      const a = await api.publishAnnounce(body)
      setAnnounces((list) => [a, ...list.filter((x) => x.id !== a.id)])
    },
    reloadAnnounces,
    reloadActivities,
  }

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>
}

export function useApp() {
  const v = useContext(Ctx)
  if (!v) throw new Error('useApp outside provider')
  return v
}
