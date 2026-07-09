import { useEffect, useState } from 'react'
import { api, type Activity } from '../lib/api'
import { useApp } from '../lib/store'

function TrackMap({ a, accent }: { a: Activity; accent: string }) {
  return (
    <div style={{ height: 120, borderRadius: 14, overflow: 'hidden', background: 'linear-gradient(160deg,#1a2330,#0f141c)', position: 'relative', marginTop: 12 }}>
      <svg viewBox="0 0 300 140" style={{ width: '100%', height: '100%' }}>
        <path d={a.route} fill="none" stroke={accent} strokeWidth="4.5" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx={a.sx} cy={a.sy} r="5" fill={accent} />
        <circle cx={a.ex} cy={a.ey} r="5" fill="#fff" stroke={accent} strokeWidth="2" />
      </svg>
      <div style={{ position: 'absolute', right: 10, bottom: 8, fontSize: 10, fontWeight: 800, color: '#9aa2ab' }}>GPS</div>
    </div>
  )
}

function weekKm(activities: Activity[]) {
  return activities.reduce((sum, a) => sum + (parseFloat(a.dist) || 0), 0)
}

export function HomeAthlete({ activities }: { activities: Activity[] }) {
  const { theme, club, announces, setScreen, user } = useApp()
  const [weekPlan, setWeekPlan] = useState('')
  const next = announces[0]
  const done = weekKm(activities)
  const doneLabel = done > 0 ? done.toFixed(1) : '0'

  useEffect(() => {
    void api.plan(1).then((p) => setWeekPlan(p.weekPlan || '')).catch(() => setWeekPlan(''))
  }, [])

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      {user?.inClub && next && (
        <div className="card" data-testid="group-preview" onClick={() => setScreen('schedule')} style={{ cursor: 'pointer' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
            <span style={{ fontSize: 11, fontWeight: 800, letterSpacing: '.6px', color: theme.accent }}>ГРУППОВЫЕ ТРЕНИРОВКИ</span>
            <span style={{ fontSize: 12, fontWeight: 700, color: theme.accent }}>Расписание →</span>
          </div>
          <div style={{ fontWeight: 700, fontSize: 14 }}>Ближайшая: {next.place}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>{next.day} · старт {next.time} · идут {next.going}</div>
        </div>
      )}
      <div className="card" style={{ display: 'flex', alignItems: 'center', gap: 10, borderRadius: 14, padding: '11px 14px' }}>
        <span style={{ width: 9, height: 9, borderRadius: '50%', background: theme.good, boxShadow: `0 0 10px ${theme.good}` }} />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 12, fontWeight: 700 }}>Активности</div>
          <div style={{ fontSize: 11, color: theme.dim }}>Клуб {club?.name || '—'}</div>
        </div>
      </div>
      <div className="card">
        <div style={{ fontSize: 13, fontWeight: 700, color: theme.dim, marginBottom: 10 }}>Эта неделя</div>
        <div style={{ fontFamily: theme.display, fontSize: 30, fontWeight: 800 }}>
          {doneLabel}
          {weekPlan ? <span style={{ fontSize: 14, color: theme.dim }}> из {weekPlan}</span> : <span style={{ fontSize: 14, color: theme.dim }}> км</span>}
        </div>
      </div>
      <div style={{ fontFamily: theme.display, fontSize: 16, fontWeight: 800, textTransform: 'uppercase' }}>Лента</div>
      {activities.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Пока нет активностей</div>
      )}
      {activities.map((a) => (
        <div key={a.id} className="card" data-testid="activity-card">
          <div style={{ fontWeight: 700 }}>{a.title}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>{a.when}</div>
          <div style={{ display: 'flex', gap: 14, marginTop: 10 }}>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.dist}</div><div style={{ fontSize: 10, color: theme.dim }}>км</div></div>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.time}</div><div style={{ fontSize: 10, color: theme.dim }}>время</div></div>
            <div><div style={{ fontFamily: theme.display, fontWeight: 800 }}>{a.pace}</div><div style={{ fontSize: 10, color: theme.dim }}>темп</div></div>
          </div>
          <TrackMap a={a} accent={theme.accent} />
        </div>
      ))}
    </div>
  )
}
