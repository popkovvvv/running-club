import type { Theme } from '../../lib/theme'

export function WeekPicker({ theme, week, weekRange, weekPlan, onChange }: {
  theme: Theme
  week: number
  weekRange: string
  weekPlan: string
  onChange: (w: number) => void
}) {
  return (
    <div className="card" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', borderRadius: 16 }}>
      <button className="btn" onClick={() => onChange(Math.max(0, week - 1))} style={{ width: 38, height: 38, background: theme.card2, color: theme.text, borderRadius: 11 }}>‹</button>
      <div style={{ textAlign: 'center' }}>
        <div style={{ fontWeight: 700 }}>{weekRange || '…'}</div>
        <div style={{ fontSize: 11, color: theme.accent }}>План {weekPlan}</div>
      </div>
      <button className="btn" onClick={() => onChange(Math.min(3, week + 1))} style={{ width: 38, height: 38, background: theme.card2, color: theme.text, borderRadius: 11 }}>›</button>
    </div>
  )
}
