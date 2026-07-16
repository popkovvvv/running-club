import { useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'
import { ScheduleCalendar, formatAnnounceDay } from '../components/schedule/ScheduleCalendar'
import { PeriodSummaryCard } from '../components/stats/PeriodSummaryCard'
import { SegmentEditor } from '../components/workout/SegmentEditor'
import { WorkoutTypePicker } from '../components/workout/WorkoutTypePicker'
import { api, type Segment, type StudentDetail } from '../lib/api'
import { buildMonthCells, MONTH_LABELS } from '../lib/monthCalendar'
import { useApp } from '../lib/store'
import { workoutTypeLabel, type WorkoutType } from '../lib/workoutTypes'

type PeriodKey = 'week' | 'month' | 'year'

export function StudentProfileScreen({ id }: { id: string }) {
  const { theme, closeOverlay, openOverlay } = useApp()
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [detail, setDetail] = useState<StudentDetail | null>(null)
  const [selectedIso, setSelectedIso] = useState<string | null>(null)
  const [period, setPeriod] = useState<PeriodKey>('week')
  const [open, setOpen] = useState(false)
  const [workoutType, setWorkoutType] = useState<WorkoutType>('easy')
  const [title, setTitle] = useState('')
  const [segs, setSegs] = useState<Segment[]>([])
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)
  const [saving, setSaving] = useState(false)

  const reload = () => void api.student(id, { year, month }).then((d) => {
    setDetail(d)
    setSelectedIso((prev) => {
      if (prev && d.planDays.some((p) => p.scheduledDate === prev)) return prev
      return d.planDays.find((p) => p.scheduledDate)?.scheduledDate || null
    })
  }).catch(() => setDetail(null))

  useEffect(() => {
    reload()
  }, [id, year, month])

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  const marked = useMemo(() => {
    const s = new Set<string>()
    detail?.planDays.forEach((p) => {
      if (p.scheduledDate) s.add(p.scheduledDate)
    })
    return s
  }, [detail])

  const completed = useMemo(() => {
    const s = new Set<string>()
    detail?.planDays.forEach((p) => {
      if (p.scheduledDate && p.status === 'completed') s.add(p.scheduledDate)
    })
    return s
  }, [detail])

  const cells = useMemo(
    () => buildMonthCells(year, month, theme, marked, completed),
    [year, month, theme, marked, completed],
  )

  const dayWorkouts = detail?.planDays.filter((p) => p.scheduledDate === selectedIso) ?? []

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

  const saveWorkout = async () => {
    if (!selectedIso) return
    setSaving(true)
    try {
      await api.createWorkout({
        kind: 'plan',
        targetUserId: id,
        workoutType,
        dayLabel: formatAnnounceDay(selectedIso),
        tag: workoutType,
        title: title || 'Индивидуальная',
        distKm: totalKm,
        weekIndex: 0,
        scheduledDate: selectedIso,
        segments: segs,
      })
      await reload()
      closeModal()
    } finally {
      setSaving(false)
    }
  }

  const modal = open && phoneEl ? createPortal(
    <div
      data-testid="coach-workout-modal"
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
          <div style={{ fontFamily: theme.display, fontWeight: 800 }}>Новое занятие</div>
          <button type="button" className="btn" onClick={closeModal} style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}>✕</button>
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
          data-testid="save-coach-workout"
          className="btn"
          disabled={saving || !selectedIso}
          onClick={() => void saveWorkout()}
          style={{ width: '100%', marginTop: 12, background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: selectedIso ? 1 : 0.5 }}
        >
          Создать занятие
        </button>
        <button className="btn" onClick={closeModal} style={{ width: '100%', marginTop: 8, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}>
          Закрыть
        </button>
      </div>
    </div>,
    phoneEl,
  ) : null

  if (!detail) {
    return <div className="card" style={{ color: theme.dim }}>Загрузка…</div>
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <button className="btn" onClick={closeOverlay} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>
      <div className="card" style={{ display: 'flex', gap: 14, alignItems: 'center' }}>
        <div style={{ width: 52, height: 52, borderRadius: '50%', background: theme.card2, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800, fontSize: 18 }}>{detail.student.init}</div>
        <div style={{ flex: 1 }}>
          <div style={{ fontWeight: 800, fontSize: 18 }}>{detail.student.name}</div>
          <div style={{ fontSize: 12, color: theme.dim, marginTop: 2 }}>Неделя</div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22, marginTop: 2 }}>
            {detail.weekKm}
            <span style={{ fontSize: 14, color: theme.dim, fontWeight: 700 }}> / {detail.weekPlanKm} км</span>
            <span style={{ fontSize: 13, color: theme.accent, marginLeft: 8 }}>{detail.comp}%</span>
          </div>
        </div>
      </div>

      {detail.summary && (
        <PeriodSummaryCard theme={theme} summary={detail.summary} period={period} onPeriodChange={setPeriod} />
      )}

      <button
        data-testid="create-student-workout"
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
        ＋ Создать занятие
      </button>

      <ScheduleCalendar
        theme={theme}
        label={`${MONTH_LABELS[month - 1]} ${year}`}
        cells={cells}
        selectedIso={selectedIso}
        onSelect={(iso) => setSelectedIso(iso)}
        onPrevMonth={() => shiftMonth(-1)}
        onNextMonth={() => shiftMonth(1)}
      />

      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase', fontSize: 14 }}>
        {selectedIso ? formatAnnounceDay(selectedIso) : 'Выберите день'}
      </div>
      {selectedIso && dayWorkouts.length === 0 && (
        <div className="card" style={{ color: theme.dim, fontSize: 13 }}>Нет тренировок в этот день</div>
      )}
      {dayWorkouts.map((d) => (
        <div
          key={d.workoutId}
          className="card"
          style={{ cursor: 'pointer', borderRadius: 16 }}
          onClick={() => openOverlay({ type: 'workout', id: d.workoutId })}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <div>
              <div style={{ fontWeight: 700 }}>{d.title}</div>
              <div style={{ fontSize: 11, color: theme.dim }}>{workoutTypeLabel(d.workoutType)}</div>
            </div>
            <div style={{ textAlign: 'right' }}>
              <div style={{ fontWeight: 800, color: theme.accent }}>
                {d.actualKm > 0 ? `${d.actualKm}/${d.plannedKm}` : d.plannedKm} км
              </div>
              <div style={{ fontSize: 14, color: d.status === 'completed' ? theme.good : theme.dim, fontWeight: 800 }}>
                {d.status === 'completed' ? '✓' : '○'}
              </div>
            </div>
          </div>
          {(d.activityWhen || d.rpe != null) && (
            <div style={{ marginTop: 8, fontSize: 12, color: theme.dim }}>
              {[d.activityWhen, d.activityPace, d.rpe != null ? `RPE ${d.rpe}` : null].filter(Boolean).join(' · ')}
            </div>
          )}
          {d.athleteReport && (
            <div style={{ marginTop: 6, fontSize: 12, color: theme.text }}>{d.athleteReport}</div>
          )}
          {d.coachComment && (
            <div style={{ marginTop: 4, fontSize: 12, color: theme.accent }}>Тренер: {d.coachComment}</div>
          )}
        </div>
      ))}
      {modal}
    </div>
  )
}
