import type { ReactNode } from 'react'
import { useApp } from '../lib/store'

type Tab = 'home' | 'plan' | 'prog' | 'races' | 'profile'

const icons: Record<Tab, ReactNode> = {
  home: (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 10.5 12 3l9 7.5" />
      <path d="M5 9.5V20h14V9.5" />
    </svg>
  ),
  plan: (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="4" width="18" height="17" rx="3" />
      <path d="M3 9h18M8 3v3M16 3v3" />
    </svg>
  ),
  prog: (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M4 20V10M10 20V4M16 20v-7M22 20H2" />
    </svg>
  ),
  races: (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="8" r="6" />
      <path d="M8.5 13.5 7 22l5-3 5 3-1.5-8.5" />
    </svg>
  ),
  profile: (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="8" r="4" />
      <path d="M5 20c0-3.5 3-6 7-6s7 2.5 7 6" />
    </svg>
  ),
}

function labels(isCoach: boolean): Record<Tab, string> {
  return {
    home: isCoach ? 'Ученики' : 'Главная',
    plan: isCoach ? 'Конструктор' : 'План',
    prog: isCoach ? 'Аналитика' : 'Прогресс',
    races: 'Старты',
    profile: isCoach ? 'Клуб' : 'Профиль',
  }
}

const tabs: Tab[] = ['home', 'plan', 'prog', 'races', 'profile']

export function BottomNav() {
  const { user, screen, setScreen } = useApp()
  const isCoach = user?.role === 'coach'
  const tabLabels = labels(!!isCoach)
  const active: Tab | null = tabs.includes(screen as Tab) ? (screen as Tab) : null

  return (
    <nav className="nav" data-testid="bottom-nav">
      {tabs.map((tab) => (
        <button
          key={tab}
          data-testid={`nav-${tab}`}
          className={active === tab ? 'active' : ''}
          onClick={() => setScreen(tab)}
        >
          {icons[tab]}
          <span>{tabLabels[tab]}</span>
        </button>
      ))}
    </nav>
  )
}
