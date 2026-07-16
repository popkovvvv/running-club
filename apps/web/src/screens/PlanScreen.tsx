import { useEffect, useState } from 'react'
import { api, type Segment } from '../lib/api'
import { useApp } from '../lib/store'
import { WORKOUT_PRESETS, type WorkoutType } from '../lib/workoutTypes'
import { SegmentEditor } from '../components/workout/SegmentEditor'
import { WeekPicker } from '../components/workout/WeekPicker'
import { WorkoutCard } from '../components/workout/WorkoutCard'
import { WorkoutTypePicker } from '../components/workout/WorkoutTypePicker'

export function PlanScreen() {
  const { theme, openOverlay } = useApp()
  const [plan, setPlan] = useState<Awaited<ReturnType<typeof api.plan>> | null>(null)
  const [week, setWeek] = useState(0)
  const [open, setOpen] = useState(false)
  const [workoutType, setWorkoutType] = useState<WorkoutType>('easy')
  const [title, setTitle] = useState('Свой кросс')
  const [dayLabel, setDayLabel] = useState('Ср')
  const [pace, setPace] = useState('7:40')
  const [duration, setDuration] = useState('~46 мин')
  const [segs, setSegs] = useState<Segment[]>(WORKOUT_PRESETS.easy)

  const reload = () => void api.plan(week).then(setPlan)

  useEffect(() => {
    reload()
  }, [week])

  useEffect(() => {
    setSegs(WORKOUT_PRESETS[workoutType])
  }, [workoutType])

  const totalKm = segs.reduce((a, s) => a + s.distKm, 0)

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <WeekPicker theme={theme} week={week} weekRange={plan?.weekRange || ''} weekPlan={plan?.weekPlan || ''} onChange={setWeek} />
      <button data-testid="create-workout" className="btn" onClick={() => setOpen(true)} style={{ border: `1.5px dashed ${theme.accent}`, background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13 }}>
        ＋ Создать свою тренировку
      </button>
      {plan?.mine?.map((w) => (
        <WorkoutCard key={w.id} theme={theme} workout={w} onClick={() => openOverlay({ type: 'workout', id: w.id })} />
      ))}
      {plan?.days?.map((d) => (
        <WorkoutCard key={d.id} theme={theme} workout={d} onClick={() => openOverlay({ type: 'workout', id: d.id })} />
      ))}
      {open && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,.55)', zIndex: 30, display: 'flex', alignItems: 'flex-end' }}>
          <div className="card" style={{ width: '100%', maxHeight: '85vh', overflowY: 'auto', borderRadius: '28px 28px 46px 46px', margin: '0 0 80px' }}>
            <div style={{ fontFamily: theme.display, fontWeight: 800, marginBottom: 12 }}>Новая тренировка</div>
            <WorkoutTypePicker theme={theme} value={workoutType} onChange={setWorkoutType} />
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginTop: 12 }}>
              <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Название" style={{ background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text }} />
              <input value={dayLabel} onChange={(e) => setDayLabel(e.target.value)} placeholder="День (Пн, Ср…)" style={{ background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text }} />
              <div style={{ display: 'flex', gap: 8 }}>
                <input value={duration} onChange={(e) => setDuration(e.target.value)} placeholder="Длительность" style={{ flex: 1, background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text }} />
                <input value={pace} onChange={(e) => setPace(e.target.value)} placeholder="Темп" style={{ flex: 1, background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text }} />
              </div>
            </div>
            <div style={{ marginTop: 12, fontSize: 12, color: theme.dim }}>Объём: {totalKm.toFixed(1)} км</div>
            <div style={{ marginTop: 8 }}>
              <SegmentEditor theme={theme} segments={segs} onChange={setSegs} />
            </div>
            <button
              data-testid="save-workout"
              className="btn"
              onClick={() => void api.createWorkout({
                kind: 'own',
                workoutType,
                dayLabel,
                tag: workoutType,
                title,
                distKm: totalKm,
                duration,
                pace,
                weekIndex: week,
                segments: segs,
              }).then(reload).then(() => setOpen(false))}
              style={{ width: '100%', marginTop: 12, background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14 }}
            >
              Добавить в план
            </button>
            <button className="btn" onClick={() => setOpen(false)} style={{ width: '100%', marginTop: 8, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}>Закрыть</button>
          </div>
        </div>
      )}
    </div>
  )
}
