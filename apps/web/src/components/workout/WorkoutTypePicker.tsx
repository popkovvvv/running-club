import { WORKOUT_TYPES, type WorkoutType } from '../../lib/workoutTypes'
import type { Theme } from '../../lib/theme'

export function WorkoutTypePicker({ theme, value, onChange }: {
  theme: Theme
  value: WorkoutType
  onChange: (v: WorkoutType) => void
}) {
  return (
    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
      {WORKOUT_TYPES.map((t) => (
        <button
          key={t.id}
          type="button"
          className="btn"
          onClick={() => onChange(t.id)}
          style={{
            padding: '8px 12px',
            borderRadius: 20,
            fontSize: 12,
            fontWeight: 700,
            background: value === t.id ? theme.accent : theme.card2,
            color: value === t.id ? theme.onAccent : theme.dim,
          }}
        >
          {t.label}
        </button>
      ))}
    </div>
  )
}
