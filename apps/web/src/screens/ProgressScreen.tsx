import { useEffect, useState } from 'react'
import { PeriodSummaryCard } from '../components/stats/PeriodSummaryCard'
import { AccentPill } from '../components/ui/AccentPill'
import { HeroPanel } from '../components/ui/HeroPanel'
import { api, type PeriodSummary } from '../lib/api'
import { useApp } from '../lib/store'

type PeriodKey = 'week' | 'month' | 'year'

export function ProgressScreen() {
  const { theme, user, openOverlay, tabEpoch } = useApp()
  const [progress, setProgress] = useState<{
    months: Array<{ m: string; km: number; tr: number; pace: string; diff: string }>
    yearKm: number
    yearTr: number
    yearStarts: number
    summary?: PeriodSummary
  } | null>(null)
  const [analytics, setAnalytics] = useState<Awaited<ReturnType<typeof api.analytics>> | null>(null)
  const [races, setRaces] = useState<Awaited<ReturnType<typeof api.races>>>([])
  const [prs, setPrs] = useState<Awaited<ReturnType<typeof api.prs>>>([])
  const [period, setPeriod] = useState<PeriodKey>('week')
  const isCoach = user?.role === 'coach'

  useEffect(() => {
    if (isCoach) {
      void api.analytics().then(setAnalytics).catch(() => setAnalytics(null))
      return
    }
    void api.progress().then(setProgress).catch(() => setProgress(null))
    void api.races().then(setRaces).catch(() => setRaces([]))
    void api.prs().then(setPrs).catch(() => setPrs([]))
  }, [isCoach, tabEpoch])

  if (isCoach) {
    return (
      <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
        <HeroPanel>
          <div className="display-title" style={{ fontSize: 22, marginBottom: 12 }}>Аналитика клуба</div>
          <div style={{ display: 'flex', gap: 20 }}>
            <div>
              <div style={{ fontFamily: theme.display, fontSize: 28, fontWeight: 800, fontStyle: 'italic' }}>{analytics?.clubKm ?? '—'}</div>
              <div style={{ fontSize: 11, color: theme.dim }}>км клуба</div>
            </div>
            <div>
              <div style={{ fontFamily: theme.display, fontSize: 28, fontWeight: 800, fontStyle: 'italic' }}>{analytics?.students?.length ?? '—'}</div>
              <div style={{ fontSize: 11, color: theme.dim }}>учеников</div>
            </div>
          </div>
        </HeroPanel>
        <div className="display-title" style={{ fontSize: 18 }}>Прогресс учеников</div>
        {(analytics?.students?.length ?? 0) === 0 && (
          <div className="card" style={{ fontSize: 13, color: theme.dim }}>Нет данных</div>
        )}
        {analytics?.students?.map((s) => (
          <div key={s.id} className="card" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', cursor: 'pointer' }} onClick={() => openOverlay({ type: 'student', id: s.id })}>
            <div><div style={{ fontWeight: 700 }}>{s.name}</div><div style={{ fontSize: 11, color: theme.dim }}>{s.km} км</div></div>
            <AccentPill>{s.comp}%</AccentPill>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      {progress?.summary && (
        <PeriodSummaryCard theme={theme} summary={progress.summary} period={period} onPeriodChange={setPeriod} />
      )}
      <HeroPanel style={{ textAlign: 'center' }}>
        <div style={{ fontSize: 12, color: theme.dim, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase' }}>Годовой итог</div>
        <div className="display-title" style={{ fontSize: 48, color: theme.accent, margin: '8px 0 4px' }}>
          {progress?.yearKm != null ? progress.yearKm : '—'}
        </div>
        <div style={{ fontSize: 13, color: theme.dim }}>
          км · {progress?.yearTr ?? '—'} трен. · {progress?.yearStarts ?? '—'} старта
        </div>
      </HeroPanel>

      <div className="display-title" style={{ fontSize: 18 }}>Старты</div>
      {races.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Нет ближайших стартов</div>
      )}
      {races.map((r) => (
        <div key={`${r.name}-${r.date}`} className="card" data-testid="race-card" style={{ display: 'flex', justifyContent: 'space-between', gap: 12 }}>
          <div>
            <div style={{ fontWeight: 700 }}>{r.name}</div>
            <div style={{ fontSize: 12, color: theme.dim, marginTop: 4 }}>{r.date} · {r.dist}</div>
            {r.goal && <div style={{ fontSize: 12, color: theme.accent, marginTop: 4 }}>{r.goal}</div>}
          </div>
          <AccentPill style={{ alignSelf: 'flex-start' }}>{r.days} дн.</AccentPill>
        </div>
      ))}

      <div className="display-title" style={{ fontSize: 18 }}>Личные рекорды</div>
      {prs.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Пока нет рекордов</div>
      )}
      {prs.map((pr) => (
        <div key={`${pr.d}-${pr.t}`} className="card" data-testid="pr-card" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontWeight: 800, fontSize: 18 }}>{pr.d}</div>
            <div style={{ fontSize: 12, color: theme.dim }}>{pr.date}</div>
          </div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 20, color: theme.accent }}>{pr.t}</div>
        </div>
      ))}

      {(progress?.months?.length ?? 0) === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Нет статистики</div>
      )}
      {progress?.months?.map((m) => (
        <div key={m.m} className="card">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <div className="display-title" style={{ fontSize: 16 }}>{m.m}</div>
            <AccentPill style={{ fontSize: 10, padding: '4px 10px' }}>{m.km} км</AccentPill>
          </div>
          <div style={{ display: 'flex', gap: 16, fontSize: 13, color: theme.dim }}>
            <span>{m.tr} трен.</span><span>{m.pace}</span><span>сложн. {m.diff}</span>
          </div>
        </div>
      ))}
    </div>
  )
}
