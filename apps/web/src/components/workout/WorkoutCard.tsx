import { workoutTypeLabel } from '../../lib/workoutTypes'
import type { Workout } from '../../lib/api'
import type { Theme } from '../../lib/theme'

export function WorkoutCard({ theme, workout, onClick }: {
  theme: Theme
  workout: Workout
  onClick?: () => void
}) {
  const statusIcon = workout.status === 'completed' ? '✓' : workout.status === 'skipped' ? '—' : '○'
  const statusColor = workout.status === 'completed' ? theme.good : theme.dim

  return (
    <div
      className="card"
      onClick={onClick}
      style={{ borderRadius: 16, cursor: onClick ? 'pointer' : undefined, borderColor: workout.kind === 'own' ? theme.accent : undefined }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{workout.dayLabel || '—'}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>{workout.tag || workoutTypeLabel(workout.workoutType)}</div>
        </div>
        <div style={{ textAlign: 'right' }}>
          <div style={{ fontWeight: 800 }}>{workout.distKm || '—'} км</div>
          <div style={{ fontSize: 11, color: theme.accent }}>{workout.duration || workout.pace}</div>
          <div style={{ fontSize: 14, color: statusColor, fontWeight: 800 }}>{statusIcon}</div>
        </div>
      </div>
      <div style={{ fontWeight: 700, marginTop: 8 }}>{workout.title}</div>
      {workout.segments && workout.segments.length > 0 && (
        <div style={{ marginTop: 8, fontSize: 11, color: theme.dim }}>
          {workout.segments.length} сегментов · {workout.segments.map((s) => s.title).join(', ')}
        </div>
      )}
    </div>
  )
}
