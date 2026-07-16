import type { Segment } from '../../lib/api'
import type { Theme } from '../../lib/theme'
import type { SegmentKind } from '../../lib/workoutTypes'
import { SegmentKindField } from './SegmentKindField'

export function SegmentEditor({ theme, segments, onChange }: {
  theme: Theme
  segments: Segment[]
  onChange: (s: Segment[]) => void
}) {
  const update = (i: number, patch: Partial<Segment>) => {
    onChange(segments.map((s, idx) => (idx === i ? { ...s, ...patch } : s)))
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      {segments.map((s, i) => (
        <div key={i} className="card" style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <SegmentKindField
              theme={theme}
              value={s.kind}
              onChange={(kind: SegmentKind) => update(i, { kind })}
            />
            <button type="button" className="btn" onClick={() => onChange(segments.filter((_, j) => j !== i))} style={{ background: theme.card2, color: theme.dim, borderRadius: 8, padding: '6px 10px' }}>✕</button>
          </div>
          <input value={s.title} onChange={(e) => update(i, { title: e.target.value })} placeholder="Название" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text, boxSizing: 'border-box' }} />
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <input type="number" inputMode="decimal" value={s.distKm || ''} onChange={(e) => update(i, { distKm: e.target.value === '' ? 0 : parseFloat(e.target.value) || 0 })} placeholder="км" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text, boxSizing: 'border-box' }} />
            <input value={s.pace} onChange={(e) => update(i, { pace: e.target.value })} placeholder="темп" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text, boxSizing: 'border-box' }} />
          </div>
        </div>
      ))}
      <button
        type="button"
        className="btn"
        onClick={() => onChange([...segments, { kind: 'work', title: '', distKm: 0, pace: '' }])}
        style={{ background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 12, border: `1px dashed ${theme.accent}` }}
      >
        ＋ Сегмент
      </button>
    </div>
  )
}
