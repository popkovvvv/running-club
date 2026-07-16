import { useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'
import { ScheduleCalendar, formatAnnounceDay } from '../components/schedule/ScheduleCalendar'
import { SegmentEditor } from '../components/workout/SegmentEditor'
import { WorkoutCard } from '../components/workout/WorkoutCard'
import { WorkoutTypePicker } from '../components/workout/WorkoutTypePicker'
import { api, type Segment, type Workout } from '../lib/api'
import { buildMonthCells, MONTH_LABELS } from '../lib/monthCalendar'
import { useApp } from '../lib/store'
import type { WorkoutType } from '../lib/workoutTypes'

export function PlanScreen() {
  const { theme, openOverlay, tabEpoch } = useApp()
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [plan, setPlan] = useState<Awaited<ReturnType<typeof api.plan>> | null>(null)
  const [selectedIso, setSelectedIso] = useState<string | null>(null)
  const [open, setOpen] = useState(false)
  const [workoutType, setWorkoutType] = useState<WorkoutType>('easy')
  const [title, setTitle] = useState('')
  const [segs, setSegs] = useState<Segment[]>([])
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)

  const reload = () => void api.plan({ year, month }).then((p) => {
    setPlan(p)
    setSelectedIso((prev) => {
      const all = [...(p.days || []), ...(p.mine || [])]
      if (prev && all.some((w) => w.scheduledDate === prev)) return prev
      return all.find((w) => w.scheduledDate)?.scheduledDate || null
    })
  })

  useEffect(() => {
    reload()
  }, [year, month, tabEpoch])

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  const allWorkouts = useMemo(() => [...(plan?.days || []), ...(plan?.mine || [])], [plan])

  const marked = useMemo(() => {
    const s = new Set<string>()
    allWorkouts.forEach((w) => {
      if (w.scheduledDate) s.add(w.scheduledDate)
    })
    return s
  }, [allWorkouts])

  const completed = useMemo(() => {
    const s = new Set<string>()
    allWorkouts.forEach((w) => {
      if (w.scheduledDate && w.status === 'completed') s.add(w.scheduledDate)
    })
    return s
  }, [allWorkouts])

  const cells = useMemo(
    () => buildMonthCells(year, month, theme, marked, completed),
    [year, month, theme, marked, completed],
  )

  const dayWorkouts = allWorkouts.filter((w) => w.scheduledDate === selectedIso)

  const shiftMonth = (delta: number) => {
    const d = new Date(year, month - 1 + delta, 1)
    setYear(d.getFullYear())
    setMonth(d.getMonth() + 1)
    setSelectedIso(null)
  }

  const resetForm = () => {
    setWorkoutType('easy')
    setTitle('')
    setSegs([])
  }

  const openModal = () => {
    if (!selectedIso) return
    resetForm()
    setOpen(true)
  }

  const closeModal = () => {
    setOpen(false)
    resetForm()
  }

  const totalKm = segs.reduce((a, s) => a + (s.distKm || 0), 0)

  const modal = open && phoneEl ? createPortal(
    <div
      data-testid="workout-modal"
      style={{
        position: 'absolute',
        inset: 0,
        background: 'rgba(0,0,0,.6)',
        zIndex: 50,
        display: 'flex',
        alignItems: 'flex-end',
      }}
      onClick={closeModal}
    >
      <div
        className="card"
        onClick={(e) => e.stopPropagation()}
        style={{
          width: '100%',
          maxHeight: '92%',
          overflowY: 'auto',
          borderRadius: '24px 24px 0 0',
          paddingBottom: 28,
          boxSizing: 'border-box',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800 }}>Новая тренировка</div>
          <button
            type="button"
            className="btn"
            data-testid="workout-modal-close"
            onClick={closeModal}
            style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}
          >
            ✕
          </button>
        </div>
        <WorkoutTypePicker theme={theme} value={workoutType} onChange={setWorkoutType} />
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginTop: 12 }}>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Название"
            style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }}
          />
        </div>
        <div style={{ marginTop: 12, fontSize: 12, color: theme.dim }}>
          Объём: {totalKm > 0 ? `${totalKm.toFixed(1)} км` : '—'}
        </div>
        <div style={{ marginTop: 8 }}>
          <SegmentEditor theme={theme} segments={segs} onChange={setSegs} />
        </div>
        <button
          data-testid="save-workout"
          className="btn"
          disabled={!selectedIso}
          onClick={() => {
            if (!selectedIso) return
            void api.createWorkout({
              kind: 'own',
              workoutType,
              dayLabel: formatAnnounceDay(selectedIso),
              tag: workoutType,
              title: title || 'Своя тренировка',
              distKm: totalKm,
              weekIndex: 0,
              scheduledDate: selectedIso,
              segments: segs,
            }).then(reload).then(closeModal)
          }}
          style={{ width: '100%', marginTop: 12, background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: selectedIso ? 1 : 0.5 }}
        >
          Добавить в план
        </button>
        <button className="btn" onClick={closeModal} style={{ width: '100%', marginTop: 8, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}>
          Закрыть
        </button>
      </div>
    </div>,
    phoneEl,
  ) : null

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      {(plan?.weekPlan || plan?.weekKm) && (
        <div className="card" data-testid="week-volume">
          <div style={{ fontSize: 11, color: theme.dim, fontWeight: 700, letterSpacing: '0.04em', textTransform: 'uppercase' }}>
            Объём недели{plan.weekRange ? ` · ${plan.weekRange}` : ''}
          </div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 28, marginTop: 6 }}>
            {plan.weekKm || '0'}
            {plan.weekPlan ? <span style={{ fontSize: 16, color: theme.dim, fontWeight: 700 }}> / {plan.weekPlan}</span> : <span style={{ fontSize: 16, color: theme.dim }}> км</span>}
          </div>
        </div>
      )}
      <ScheduleCalendar
        theme={theme}
        label={`${MONTH_LABELS[month - 1]} ${year}`}
        cells={cells}
        selectedIso={selectedIso}
        onSelect={(iso) => setSelectedIso(iso)}
        onPrevMonth={() => shiftMonth(-1)}
        onNextMonth={() => shiftMonth(1)}
      />
      <button
        data-testid="create-workout"
        className="btn"
        disabled={!selectedIso}
        onClick={openModal}
        style={{
          border: `1.5px dashed ${theme.accent}`,
          background: theme.accentSoft,
          color: theme.accent,
          borderRadius: 14,
          padding: 13,
          opacity: selectedIso ? 1 : 0.5,
        }}
      >
        ＋ Создать свою тренировку
      </button>
      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase', fontSize: 14 }}>
        {selectedIso ? formatAnnounceDay(selectedIso) : 'Выберите день'}
      </div>
      {selectedIso && dayWorkouts.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Нет тренировок в этот день</div>
      )}
      {dayWorkouts.map((w: Workout) => (
        <WorkoutCard key={w.id} theme={theme} workout={w} onClick={() => openOverlay({ type: 'workout', id: w.id })} />
      ))}
      {modal}
    </div>
  )
}
