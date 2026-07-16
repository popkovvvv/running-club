import { useEffect, useState } from 'react'
import { api, type Segment, type Student, type Workout } from '../lib/api'
import { useApp } from '../lib/store'
import { type WorkoutType } from '../lib/workoutTypes'
import { DateField } from '../components/workout/DateField'
import { SegmentEditor } from '../components/workout/SegmentEditor'
import { StudentField } from '../components/workout/StudentField'
import { WeekPicker } from '../components/workout/WeekPicker'
import { WorkoutCard } from '../components/workout/WorkoutCard'
import { WorkoutTypePicker } from '../components/workout/WorkoutTypePicker'

type CoachTab = 'template' | 'students' | 'assign'

export function CoachPlanScreen() {
  const { theme, openOverlay, tabEpoch } = useApp()
  const [tab, setTab] = useState<CoachTab>('template')
  const [week, setWeek] = useState(0)
  const [weekMeta, setWeekMeta] = useState({ rangeLabel: '', planLabel: '' })
  const [template, setTemplate] = useState<Workout[]>([])
  const [students, setStudents] = useState<Student[]>([])
  const [selectedStudent, setSelectedStudent] = useState('')
  const [workoutType, setWorkoutType] = useState<WorkoutType>('interval')
  const [title, setTitle] = useState('')
  const [dayLabel, setDayLabel] = useState('')
  const [scheduledDate, setScheduledDate] = useState<string | null>(null)
  const [segs, setSegs] = useState<Segment[]>([])
  const [draftWorkouts, setDraftWorkouts] = useState<Workout[]>([])

  const load = () => {
    void api.planWeeks().then((weeks) => {
      const w = weeks.find((x) => x.weekIndex === week)
      if (w) setWeekMeta({ rangeLabel: w.rangeLabel, planLabel: w.planLabel })
    })
    void api.planTemplate(week).then((t) => {
      setTemplate(t.workouts)
      setDraftWorkouts(t.workouts)
    })
    void api.students().then(setStudents).catch(() => setStudents([]))
  }

  useEffect(() => {
    load()
  }, [week, tabEpoch])

  const totalKm = segs.reduce((a, s) => a + (s.distKm || 0), 0)

  const addToTemplate = () => {
    const item: Workout = {
      id: `draft-${Date.now()}`,
      kind: 'plan',
      workoutType,
      dayLabel,
      tag: workoutType,
      title,
      distKm: totalKm,
      hr: '',
      weekIndex: week,
      status: 'planned',
      scheduledDate: scheduledDate || undefined,
      segments: segs,
    }
    setDraftWorkouts((list) => [...list, item])
  }

  const saveTemplate = async () => {
    await api.savePlanTemplate(week, draftWorkouts.map((w) => ({
      kind: 'plan',
      workoutType: w.workoutType,
      dayLabel: w.dayLabel,
      tag: w.tag,
      title: w.title,
      distKm: w.distKm,
      weekIndex: week,
      scheduledDate: w.scheduledDate,
      segments: w.segments || [],
    })))
    load()
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div style={{ display: 'flex', gap: 6 }}>
        {(['template', 'students', 'assign'] as CoachTab[]).map((t) => (
          <button key={t} className="btn" onClick={() => setTab(t)} style={{ flex: 1, padding: 10, borderRadius: 12, fontSize: 11, fontWeight: 700, background: tab === t ? theme.accent : theme.card2, color: tab === t ? theme.onAccent : theme.dim }}>
            {t === 'template' ? 'Шаблон' : t === 'students' ? 'Ученики' : 'Назначить'}
          </button>
        ))}
      </div>
      <WeekPicker theme={theme} week={week} weekRange={weekMeta.rangeLabel} weekPlan={weekMeta.planLabel} onChange={setWeek} />

      {tab === 'template' && (
        <>
          <div style={{ display: 'flex', gap: 8 }}>
            <input value={weekMeta.rangeLabel} onChange={(e) => setWeekMeta((m) => ({ ...m, rangeLabel: e.target.value }))} placeholder="Диапазон дат" style={{ flex: 1, background: theme.card2, border: 'none', borderRadius: 10, padding: 10, color: theme.text }} />
            <input value={weekMeta.planLabel} onChange={(e) => setWeekMeta((m) => ({ ...m, planLabel: e.target.value }))} placeholder="Объём" style={{ width: 90, background: theme.card2, border: 'none', borderRadius: 10, padding: 10, color: theme.text }} />
            <button className="btn" onClick={() => void api.upsertPlanWeek(week, weekMeta).then(load)} style={{ background: theme.card2, color: theme.accent, borderRadius: 10, padding: '8px 12px' }}>OK</button>
          </div>
          <WorkoutTypePicker theme={theme} value={workoutType} onChange={setWorkoutType} />
          <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Название" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }} />
          <DateField
            theme={theme}
            value={scheduledDate}
            dayLabel={dayLabel}
            onChange={(iso, label) => {
              setScheduledDate(iso)
              setDayLabel(label)
            }}
          />
          <div style={{ fontSize: 12, color: theme.dim }}>Объём: {totalKm > 0 ? `${totalKm.toFixed(1)} км` : '—'}</div>
          <SegmentEditor theme={theme} segments={segs} onChange={setSegs} />
          <button className="btn" onClick={addToTemplate} style={{ background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 12 }}>＋ В шаблон недели</button>
          {draftWorkouts.map((w) => (
            <WorkoutCard key={w.id} theme={theme} workout={w} />
          ))}
          <button className="btn" onClick={() => void saveTemplate()} style={{ background: theme.card2, color: theme.text, borderRadius: 14, padding: 14 }}>Сохранить шаблон</button>
          <button className="btn" onClick={() => void api.publishPlan(week).then(load)} style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 15 }}>Опубликовать клубу</button>
          {template.length > 0 && <div style={{ fontSize: 11, color: theme.dim }}>Сохранено: {template.length} тренировок</div>}
        </>
      )}

      {tab === 'students' && students.map((s) => (
        <div key={s.id} className="card" style={{ display: 'flex', justifyContent: 'space-between', cursor: 'pointer' }} onClick={() => openOverlay({ type: 'student', id: s.id })}>
          <div><div style={{ fontWeight: 700 }}>{s.name}</div><div style={{ fontSize: 11, color: theme.dim }}>{s.km} км</div></div>
          <div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{s.comp}%</div>
        </div>
      ))}

      {tab === 'assign' && (
        <>
          <StudentField theme={theme} students={students} value={selectedStudent} onChange={setSelectedStudent} />
          <WorkoutTypePicker theme={theme} value={workoutType} onChange={setWorkoutType} />
          <DateField
            theme={theme}
            value={scheduledDate}
            dayLabel={dayLabel}
            onChange={(iso, label) => {
              setScheduledDate(iso)
              setDayLabel(label)
            }}
          />
          <SegmentEditor theme={theme} segments={segs} onChange={setSegs} />
          <button
            className="btn"
            disabled={!selectedStudent}
            onClick={() => void api.createWorkout({
              kind: 'plan',
              targetUserId: selectedStudent,
              workoutType,
              dayLabel,
              title: title || 'Индивидуальная',
              distKm: totalKm,
              weekIndex: week,
              scheduledDate: scheduledDate || undefined,
              segments: segs,
            })}
            style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 15, opacity: selectedStudent ? 1 : 0.5 }}
          >
            Назначить ученику
          </button>
        </>
      )}
    </div>
  )
}
