import type { PeriodStats, PeriodSummary } from '../../lib/api'
import type { Theme } from '../../lib/theme'

type PeriodKey = keyof PeriodSummary

const LABELS: Record<PeriodKey, string> = {
  week: 'Неделя',
  month: 'Месяц',
  year: 'Год',
}

export function PeriodSummaryCard({
  theme,
  summary,
  period,
  onPeriodChange,
}: {
  theme: Theme
  summary: PeriodSummary
  period: PeriodKey
  onPeriodChange: (p: PeriodKey) => void
}) {
  const stats: PeriodStats = summary[period]
  return (
    <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div style={{ display: 'flex', gap: 8 }}>
        {(Object.keys(LABELS) as PeriodKey[]).map((key) => (
          <button
            key={key}
            type="button"
            className="btn"
            onClick={() => onPeriodChange(key)}
            style={{
              flex: 1,
              padding: 8,
              borderRadius: 10,
              fontSize: 11,
              background: period === key ? theme.accent : theme.card2,
              color: period === key ? theme.onAccent : theme.dim,
            }}
          >
            {LABELS[key]}
          </button>
        ))}
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
        <div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22 }}>{stats.km}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>км</div>
        </div>
        <div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22 }}>{stats.workouts}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>тренировок</div>
        </div>
        <div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22 }}>{stats.planCompletion}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>% плана</div>
        </div>
        <div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22 }}>{stats.avgPace}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>ср. темп</div>
        </div>
      </div>
    </div>
  )
}
