import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { useApp } from '../lib/store'

export function ProgressScreen() {
  const { theme, user } = useApp()
  const [progress, setProgress] = useState<Awaited<ReturnType<typeof api.progress>> | null>(null)
  const [analytics, setAnalytics] = useState<Awaited<ReturnType<typeof api.analytics>> | null>(null)
  const isCoach = user?.role === 'coach'

  useEffect(() => {
    if (isCoach) void api.analytics().then(setAnalytics).catch(() => null)
    else void api.progress().then(setProgress).catch(() => null)
  }, [isCoach])

  if (isCoach) {
    return (
      <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
        <div style={{ display: 'flex', gap: 10 }}>
          <div className="card" style={{ flex: 1 }}><div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800 }}>{analytics?.clubKm ?? '—'}</div><div style={{ fontSize: 11, color: theme.dim }}>км клуба</div></div>
          <div className="card" style={{ flex: 1 }}><div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800 }}>{analytics?.attendance ?? '—'}%</div><div style={{ fontSize: 11, color: theme.dim }}>явка</div></div>
        </div>
        <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase' }}>Прогресс учеников</div>
        {analytics?.students?.map((s) => (
          <div key={s.id} className="card" style={{ display: 'flex', justifyContent: 'space-between' }}>
            <div><div style={{ fontWeight: 700 }}>{s.name}</div><div style={{ fontSize: 11, color: theme.dim }}>{s.km} км</div></div>
            <div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{s.comp}%</div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <div className="card" style={{ textAlign: 'center' }}>
        <div style={{ fontSize: 12, color: theme.dim }}>Годовой итог</div>
        <div style={{ fontFamily: theme.display, fontSize: 42, fontWeight: 800, color: theme.accent }}>{progress?.yearKm ?? 612}</div>
        <div style={{ fontSize: 13, color: theme.dim }}>км · {progress?.yearTr ?? 78} трен. · {progress?.yearStarts ?? 3} старта</div>
      </div>
      {progress?.months?.map((m) => (
        <div key={m.m} className="card">
          <div style={{ fontFamily: theme.display, fontWeight: 800 }}>{m.m}</div>
          <div style={{ display: 'flex', gap: 16, marginTop: 8, fontSize: 13 }}>
            <span>{m.km} км</span><span>{m.tr} трен.</span><span>{m.pace}</span><span>сложн. {m.diff}</span>
          </div>
        </div>
      ))}
    </div>
  )
}
