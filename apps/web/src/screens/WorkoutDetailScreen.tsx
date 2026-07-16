import { useEffect, useState } from 'react'
import { api, type Workout } from '../lib/api'
import { useApp } from '../lib/store'
import { segmentKindLabel, workoutTypeLabel } from '../lib/workoutTypes'

export function WorkoutDetailScreen({ id }: { id: string }) {
  const { theme, closeOverlay, openOverlay } = useApp()
  const [workout, setWorkout] = useState<Workout | null>(null)

  useEffect(() => {
    void api.getWorkout(id).then(setWorkout).catch(() => setWorkout(null))
  }, [id])

  if (!workout) {
    return <div className="card" style={{ color: theme.dim }}>Загрузка…</div>
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <button className="btn" onClick={closeOverlay} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>
      <div className="card">
        <div style={{ fontSize: 11, color: theme.accent, fontWeight: 700 }}>{workoutTypeLabel(workout.workoutType)} · {workout.dayLabel}</div>
        <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22, marginTop: 6 }}>{workout.title}</div>
        {workout.description && <div style={{ fontSize: 13, color: theme.dim, marginTop: 8 }}>{workout.description}</div>}
        <div style={{ display: 'flex', gap: 16, marginTop: 12 }}>
          <div><div style={{ fontWeight: 800 }}>{workout.distKm} км</div><div style={{ fontSize: 10, color: theme.dim }}>дистанция</div></div>
          <div><div style={{ fontWeight: 800 }}>{workout.pace || '—'}</div><div style={{ fontSize: 10, color: theme.dim }}>темп</div></div>
          <div><div style={{ fontWeight: 800 }}>{workout.status === 'completed' ? '✓' : '○'}</div><div style={{ fontSize: 10, color: theme.dim }}>статус</div></div>
        </div>
      </div>
      {workout.segments?.map((s, i) => (
        <div key={s.id || i} className="card" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <div><div style={{ fontSize: 11, color: theme.accent }}>{segmentKindLabel(s.kind)}</div><div style={{ fontWeight: 700 }}>{s.title}</div></div>
          <div style={{ textAlign: 'right' }}><div style={{ fontWeight: 800 }}>{s.distKm} км</div><div style={{ fontSize: 11, color: theme.dim }}>{s.pace}</div></div>
        </div>
      ))}
      {workout.completedActivityId && (
        <button className="btn" onClick={() => openOverlay({ type: 'activity', id: workout.completedActivityId! })} style={{ background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 14 }}>
          Открыть активность
        </button>
      )}
      {workout.status !== 'completed' && (
        <button
          className="btn"
          onClick={() => void api.updateWorkout(id, { status: 'completed' }).then(setWorkout)}
          style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14 }}
        >
          Отметить выполненной
        </button>
      )}
      {workout.kind === 'own' && (
        <button className="btn" onClick={() => void api.deleteWorkout(id).then(closeOverlay)} style={{ background: theme.card2, color: theme.accent, borderRadius: 14, padding: 12 }}>
          Удалить
        </button>
      )}
    </div>
  )
}
