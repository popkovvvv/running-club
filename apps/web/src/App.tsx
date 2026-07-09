import { useEffect, useState, type CSSProperties } from 'react'
import { AuthGate } from './components/AuthGate'
import { BottomNav } from './components/BottomNav'
import { CreateClubGate } from './components/CreateClubGate'
import { RoleEntry } from './components/RoleEntry'
import { api, type Activity } from './lib/api'
import { useApp } from './lib/store'
import { themeToCssVars } from './lib/theme'
import { HomeAthlete } from './screens/HomeAthlete'
import { HomeCoach } from './screens/HomeCoach'
import { PlanScreen } from './screens/PlanScreen'
import { ProfileAthlete, ProfileCoach } from './screens/ProfileScreen'
import { ProgressScreen } from './screens/ProgressScreen'
import { RacesScreen } from './screens/RacesScreen'
import { ScheduleScreen } from './screens/ScheduleScreen'

type AuthRole = 'athlete' | 'coach'

function initials(name?: string, isCoach?: boolean) {
  if (!name?.trim()) return isCoach ? 'ТР' : '—'
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return parts[0].slice(0, 2).toUpperCase()
}

export default function App() {
  const { user, loading, theme, screen, setScreen } = useApp()
  const [authRole, setAuthRole] = useState<AuthRole | null>(null)
  const [activities, setActivities] = useState<Activity[]>([])
  const isCoach = user?.role === 'coach'

  useEffect(() => {
    if (user && !isCoach) void api.activities().then(setActivities).catch(() => setActivities([]))
  }, [user, isCoach])

  const vars = themeToCssVars(theme) as CSSProperties

  if (loading) {
    return <div className="phone" style={vars}><div style={{ margin: 'auto', color: theme.dim }}>Загрузка…</div></div>
  }

  if (!user) {
    return (
      <div className="phone" style={vars} data-testid="phone-app">
        {!authRole
          ? <RoleEntry onPick={setAuthRole} />
          : <AuthGate role={authRole} onBack={() => setAuthRole(null)} />}
      </div>
    )
  }

  if (user.needsClub) {
    return (
      <div className="phone" style={vars} data-testid="phone-app">
        <CreateClubGate />
      </div>
    )
  }

  return (
    <div className="phone" style={vars} data-testid="phone-app">
      <div style={{ display: 'flex', justifyContent: 'space-between', padding: '14px 26px 6px', fontSize: 14, fontWeight: 700 }}>
        <span>9:41</span>
        <span style={{ width: 22, height: 11, border: `1.5px solid ${theme.text}`, borderRadius: 3, display: 'inline-block' }} />
      </div>
      <div style={{ padding: '8px 18px 12px', display: 'flex', flexDirection: 'column', gap: 12 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontSize: 10, color: theme.dim, fontWeight: 700, letterSpacing: 1 }}>КЛУБ «{theme.name}»</div>
            <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800 }}>{isCoach ? 'Кабинет тренера' : 'Личный кабинет'}</div>
          </div>
          <div
            style={{ width: 42, height: 42, borderRadius: '50%', background: theme.card2, border: `1.5px solid ${theme.accent}`, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800 }}
            onClick={() => setScreen('profile')}
          >
            {initials(user.name, isCoach)}
          </div>
        </div>
      </div>
      <div className="scrl" style={{ flex: 1, overflowY: 'auto', padding: '2px 16px 96px' }}>
        {screen === 'home' && (isCoach ? <HomeCoach /> : <HomeAthlete activities={activities} />)}
        {screen === 'schedule' && <ScheduleScreen />}
        {screen === 'plan' && <PlanScreen />}
        {screen === 'prog' && <ProgressScreen />}
        {screen === 'races' && <RacesScreen />}
        {screen === 'profile' && (isCoach ? <ProfileCoach /> : <ProfileAthlete />)}
      </div>
      <BottomNav />
    </div>
  )
}
