import { useState, type CSSProperties } from 'react'
import { AuthGate } from './components/AuthGate'
import { BottomNav } from './components/BottomNav'
import { CreateClubGate } from './components/CreateClubGate'
import { RoleEntry } from './components/RoleEntry'
import { useApp } from './lib/store'
import { themeToCssVars } from './lib/theme'
import { ActivityDetailScreen } from './screens/ActivityDetailScreen'
import { CoachPlanScreen } from './screens/CoachPlanScreen'
import { HomeAthlete } from './screens/HomeAthlete'
import { HomeCoach } from './screens/HomeCoach'
import { PlanScreen } from './screens/PlanScreen'
import { ProfileAthlete, ProfileCoach } from './screens/ProfileScreen'
import { ProgressScreen } from './screens/ProgressScreen'
import { ScheduleScreen } from './screens/ScheduleScreen'
import { StudentProfileScreen } from './screens/StudentProfileScreen'
import { WorkoutDetailScreen } from './screens/WorkoutDetailScreen'

type AuthRole = 'athlete' | 'coach'

function initials(name?: string, isCoach?: boolean) {
  if (!name?.trim()) return isCoach ? 'ТР' : '—'
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return parts[0].slice(0, 2).toUpperCase()
}

export default function App() {
  const { user, loading, theme, screen, overlay, setScreen } = useApp()
  const [authRole, setAuthRole] = useState<AuthRole | null>(null)
  const isCoach = user?.role === 'coach'

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

  const mainContent = overlay ? (
    <>
      {overlay.type === 'activity' && <ActivityDetailScreen id={overlay.id} />}
      {overlay.type === 'workout' && <WorkoutDetailScreen id={overlay.id} />}
      {overlay.type === 'student' && <StudentProfileScreen id={overlay.id} />}
    </>
  ) : (
    <>
      {screen === 'home' && (isCoach ? <HomeCoach /> : <HomeAthlete />)}
      {screen === 'schedule' && <ScheduleScreen />}
      {screen === 'plan' && (isCoach ? <CoachPlanScreen /> : <PlanScreen />)}
      {screen === 'prog' && <ProgressScreen />}
      {screen === 'profile' && (isCoach ? <ProfileCoach /> : <ProfileAthlete />)}
    </>
  )

  return (
    <div className="phone" style={vars} data-testid="phone-app">
      <div className="phone-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ minWidth: 0, flex: 1, paddingRight: 12 }}>
          <div className="display-title" style={{ fontSize: 22, color: theme.text }}>{theme.name}</div>
          <div style={{ fontSize: 12, color: theme.dim, fontWeight: 700, marginTop: 2 }}>{isCoach ? 'Кабинет тренера' : 'Личный кабинет'}</div>
        </div>
        <div
          style={{ width: 42, height: 42, flex: 'none', borderRadius: '50%', background: theme.card2, border: `1.5px solid ${theme.accent}`, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800 }}
          onClick={() => setScreen('profile')}
        >
          {initials(user.name, isCoach)}
        </div>
      </div>
      <div className="scrl phone-scroll">
        {mainContent}
      </div>
      {!overlay && <BottomNav />}
    </div>
  )
}
