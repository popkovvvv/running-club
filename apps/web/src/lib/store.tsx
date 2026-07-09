import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { api, type Announce, type Club, type User } from './api'
import { pulseTheme, withAccent, type Theme } from './theme'

type Screen = 'home' | 'schedule' | 'plan' | 'prog' | 'races' | 'profile'

type AppState = {
  user: User | null
  token: string | null
  theme: Theme
  club: Club | null
  screen: Screen
  announces: Announce[]
  loading: boolean
  setScreen: (s: Screen) => void
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
  publishAnnounce: (body: { place: string; day: string; time: string; group: string; note: string }) => Promise<void>
}

const Ctx = createContext<AppState | null>(null)

export function AppProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [club, setClub] = useState<Club | null>(null)
  const [screen, setScreen] = useState<Screen>('home')
  const [announces, setAnnounces] = useState<Announce[]>([])
  const [loading, setLoading] = useState(true)

  const theme = useMemo(() => withAccent(pulseTheme, club?.accentHex || pulseTheme.accent), [club?.accentHex])

  const refresh = useCallback(async () => {
    if (!localStorage.getItem('token')) {
      setUser(null)
      setClub(null)
      setAnnounces([])
      setLoading(false)
      return
    }
    try {
      const me = await api.me()
      setUser(me)
      if (me.needsClub) {
        setClub(null)
        setAnnounces([])
      } else if (me.inClub || me.role === 'coach') {
        const c = await api.club()
        setClub(c)
        const a = await api.announces()
        setAnnounces(a)
      } else {
        setClub(null)
        setAnnounces([])
      }
    } catch {
      localStorage.removeItem('token')
      setToken(null)
      setUser(null)
    } finally {
      setLoading(false)
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
    announces,
    loading,
    setScreen,
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
      setAnnounces((list) => [a, ...list])
    },
  }

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>
}

export function useApp() {
  const v = useContext(Ctx)
  if (!v) throw new Error('useApp outside provider')
  return v
}
