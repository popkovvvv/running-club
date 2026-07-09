import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { useApp } from '../lib/store'

export function ScheduleScreen() {
  const { theme, announces, signupAnnounce, unsignupAnnounce, demoRole, publishAnnounce } = useApp()
  const [cells, setCells] = useState<Array<{ key: string; n: number; blank: boolean; bg: string; fg: string; dot: string }>>([])
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    void api.calendar().then((r) => setCells(r.cells))
  }, [])

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <div className="card">
        <div style={{ fontFamily: theme.display, fontWeight: 800, marginBottom: 12 }}>Июль 2026</div>
        <div data-testid="calendar-grid" style={{ display: 'grid', gridTemplateColumns: 'repeat(7,1fr)', gap: 6 }}>
          {cells.map((c) => (
            <div
              key={c.key}
              data-testid={c.blank ? 'cal-blank' : 'cal-day'}
              data-blank={c.blank ? '1' : '0'}
              style={{ aspectRatio: '1', borderRadius: 9, background: c.bg, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', position: 'relative' }}
            >
              <span style={{ fontSize: 12, fontWeight: 700, color: c.fg }}>{c.n || ''}</span>
              {!c.blank && <span data-testid="cal-dot" style={{ width: 5, height: 5, borderRadius: '50%', background: c.dot, position: 'absolute', bottom: 5 }} />}
            </div>
          ))}
        </div>
      </div>

      {demoRole === 'coach' && (
        <button className="btn" onClick={() => setShowForm((v) => !v)} style={{ border: `1px solid ${theme.accent}`, background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13 }}>
          ＋ Опубликовать анонс
        </button>
      )}
      {showForm && (
        <button
          data-testid="publish-announce"
          className="btn"
          onClick={() => void publishAnnounce({ place: 'Парк Победы', day: 'Сб, 25 июля', time: '10:00', group: 'Основная группа', note: 'Лёгкий кросс.' }).then(() => setShowForm(false))}
          style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14 }}
        >
          Опубликовать для учеников
        </button>
      )}

      <span style={{ fontFamily: theme.display, fontSize: 16, fontWeight: 800, textTransform: 'uppercase' }}>Ближайшие занятия</span>
      {announces.map((an) => (
        <div key={an.id} className="card" data-testid="announce-card" style={{ position: 'relative', overflow: 'hidden' }}>
          <div style={{ position: 'absolute', left: 0, top: 0, bottom: 0, width: 4, background: theme.accent }} />
          <div style={{ fontSize: 10, fontWeight: 800, color: theme.accent }}>ГРУППОВАЯ · {an.group}</div>
          <div style={{ fontFamily: theme.display, fontSize: 19, fontWeight: 800, margin: '6px 0' }}>{an.place}</div>
          <div style={{ fontSize: 13, color: theme.dim, marginBottom: 14 }}>{an.day} · {an.time} · идут {an.going}</div>
          {demoRole === 'athlete' && (
            <button
              data-testid="schedule-cta"
              className="btn"
              onClick={() => void (an.signedUp ? unsignupAnnounce(an.id) : signupAnnounce(an.id))}
              style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 11, padding: 12, fontSize: 13 }}
            >
              {an.scheduleCta}
            </button>
          )}
        </div>
      ))}
    </div>
  )
}
