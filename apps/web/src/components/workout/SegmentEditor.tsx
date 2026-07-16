import { SEGMENT_KINDS, type SegmentKind } from '../../lib/workoutTypes'
import type { Theme } from '../../lib/theme'
import type { Segment } from '../../lib/api'

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
            <select
              value={s.kind}
              onChange={(e) => update(i, { kind: e.target.value as SegmentKind })}
              style={{ background: theme.card2, color: theme.text, border: 'none', borderRadius: 8, padding: '6px 8px', fontSize: 12 }}
            >
              {SEGMENT_KINDS.map((k) => (
                <option key={k.id} value={k.id}>{k.label}</option>
              ))}
            </select>
            <button type="button" className="btn" onClick={() => onChange(segments.filter((_, j) => j !== i))} style={{ background: theme.card2, color: theme.dim, borderRadius: 8, padding: '6px 10px' }}>✕</button>
          </div>
          <input value={s.title} onChange={(e) => update(i, { title: e.target.value })} placeholder="Название" style={{ background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text }} />
          <div style={{ display: 'flex', gap: 8 }}>
            <input type="number" value={s.distKm} onChange={(e) => update(i, { distKm: parseFloat(e.target.value) || 0 })} placeholder="км" style={{ flex: 1, background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text }} />
            <input value={s.pace} onChange={(e) => update(i, { pace: e.target.value })} placeholder="темп" style={{ flex: 1, background: theme.card2, border: 'none', borderRadius: 8, padding: 10, color: theme.text }} />
          </div>
        </div>
      ))}
      <button
        type="button"
        className="btn"
        onClick={() => onChange([...segments, { kind: 'work', title: 'Сегмент', distKm: 1, pace: '7:00' }])}
        style={{ background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 12, border: `1px dashed ${theme.accent}` }}
      >
        ＋ Сегмент
      </button>
    </div>
  )
}
