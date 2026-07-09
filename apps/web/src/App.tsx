import { useEffect, useState, type CSSProperties } from 'react'
import { AuthGate } from './components/AuthGate'
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

export default function App() {
  const { user, loading, theme, screen, setScreen, demoRole, setDemoRole } = useApp()
  const [activities, setActivities] = useState<Activity[]>([])
  const isCoach = demoRole === 'coach'

  useEffect(() => {
    if (user && !isCoach) void api.activities().then(setActivities).catch(() => setActivities([]))
  }, [user, isCoach])

  const vars = themeToCssVars(theme) as CSSProperties

  if (loading) {
    return <div className="phone" style={vars}><div style={{ margin: 'auto', color: theme.dim }}>Загрузка…</div></div>
  }

  return (
    <div className="phone" style={vars} data-testid="phone-app">
      {!user && <AuthGate />}
      {user && (
        <>
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
              <div style={{ width: 42, height: 42, borderRadius: '50%', background: theme.card2, border: `1.5px solid ${theme.accent}`, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800 }}>{isCoach ? 'ТР' : 'НП'}</div>
            </div>
            <div style={{ display: 'flex', background: theme.card, border: `1px solid ${theme.line}`, borderRadius: 14, padding: 4, gap: 4 }}>
              <button className="btn" data-testid="role-athlete" onClick={() => void setDemoRole('athlete')} style={{ flex: 1, padding: '9px 0', borderRadius: 10, background: !isCoach ? theme.accent : 'transparent', color: !isCoach ? theme.onAccent : theme.dim }}>Спортсмен</button>
              <button className="btn" data-testid="role-coach" onClick={() => void setDemoRole('coach')} style={{ flex: 1, padding: '9px 0', borderRadius: 10, background: isCoach ? theme.accent : 'transparent', color: isCoach ? theme.onAccent : theme.dim }}>Тренер</button>
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
          <nav className="nav" data-testid="bottom-nav">
            <button data-testid="nav-home" className={screen === 'home' ? 'active' : ''} onClick={() => setScreen('home')}>{isCoach ? 'Ученики' : 'Главная'}</button>
            <button data-testid="nav-schedule" className={screen === 'schedule' ? 'active' : ''} onClick={() => setScreen('schedule')}>Расписание</button>
            <button data-testid="nav-plan" className={screen === 'plan' ? 'active' : ''} onClick={() => setScreen('plan')}>{isCoach ? 'Конструктор' : 'План'}</button>
            <button data-testid="nav-prog" className={screen === 'prog' ? 'active' : ''} onClick={() => setScreen('prog')}>{isCoach ? 'Аналитика' : 'Прогресс'}</button>
            <button data-testid="nav-races" className={screen === 'races' ? 'active' : ''} onClick={() => setScreen('races')}>Старты</button>
            <button data-testid="nav-profile" className={screen === 'profile' ? 'active' : ''} onClick={() => setScreen('profile')}>{isCoach ? 'Клуб' : 'Профиль'}</button>
          </nav>
        </>
      )}
    </div>
  )
}
