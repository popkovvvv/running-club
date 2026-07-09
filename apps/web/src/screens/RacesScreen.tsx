import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { useApp } from '../lib/store'

export function RacesScreen() {
  const { theme, demoRole } = useApp()
  const [prs, setPrs] = useState<Array<{ d: string; t: string; date: string }>>([])
  const [races, setRaces] = useState<Array<{ days: string; name: string; date: string; dist: string; goal: string }>>([])

  useEffect(() => {
    void api.prs().then(setPrs).catch(() => [])
    void api.races().then(setRaces).catch(() => [])
  }, [])

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase' }}>Личные рекорды</div>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10 }}>
        {prs.map((p) => (
          <div key={p.d} className="card" style={{ borderRadius: 16 }}>
            <div style={{ fontSize: 12, color: theme.dim }}>{p.d}</div>
            <div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800 }}>{p.t}</div>
            <div style={{ fontSize: 10, color: theme.dim }}>{p.date}</div>
          </div>
        ))}
      </div>
      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase' }}>{demoRole === 'coach' ? 'Старты клуба' : 'Календарь стартов'}</div>
      {races.map((r) => (
        <div key={r.name} className="card" style={{ display: 'flex', gap: 14, alignItems: 'center', borderRadius: 16 }}>
          <div style={{ width: 52, textAlign: 'center', background: theme.card2, borderRadius: 12, padding: 8 }}>
            <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800, color: theme.accent }}>{r.days}</div>
            <div style={{ fontSize: 9, color: theme.dim }}>дней</div>
          </div>
          <div style={{ flex: 1 }}><div style={{ fontWeight: 700 }}>{r.name}</div><div style={{ fontSize: 11, color: theme.dim }}>{r.date} · {r.dist}</div></div>
          <div style={{ fontSize: 11, fontWeight: 700, color: theme.accent, background: theme.accentSoft, padding: '5px 9px', borderRadius: 20 }}>{r.goal}</div>
        </div>
      ))}
    </div>
  )
}
