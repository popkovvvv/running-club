import { useEffect, useState } from 'react'
import { WorkoutFactPanel } from '../components/workout/WorkoutFactPanel'
import { api, type Workout } from '../lib/api'
import { useApp } from '../lib/store'
import { segmentKindLabel, workoutTypeLabel } from '../lib/workoutTypes'

type Flash = { kind: 'ok' | 'err'; text: string } | null

export function WorkoutDetailScreen({ id }: { id: string }) {
  const { theme, closeOverlay, reloadActivities, user } = useApp()
  const [workout, setWorkout] = useState<Workout | null>(null)
  const [rpe, setRpe] = useState(5)
  const [athleteReport, setAthleteReport] = useState('')
  const [coachComment, setCoachComment] = useState('')
  const [saving, setSaving] = useState(false)
  const [flash, setFlash] = useState<Flash>(null)

  useEffect(() => {
    void api.getWorkout(id).then((w) => {
      setWorkout(w)
      setRpe(w.rpe ?? 5)
      setAthleteReport(w.athleteReport || '')
      setCoachComment(w.coachComment || '')
    }).catch(() => setWorkout(null))
  }, [id])

  useEffect(() => {
    if (!flash || flash.kind !== 'ok') return
    const t = window.setTimeout(() => setFlash(null), 2500)
    return () => window.clearTimeout(t)
  }, [flash])

  if (!workout) {
    return <div className="card" style={{ color: theme.dim }}>Загрузка…</div>
  }

  const isCoach = user?.role === 'coach'
  const canComplete = workout.status !== 'completed' && !isCoach
  const canEditReport = !isCoach

  const showOk = (text: string) => setFlash({ kind: 'ok', text })
  const showErr = (text: string) => setFlash({ kind: 'err', text })

  const complete = async () => {
    setSaving(true)
    setFlash(null)
    try {
      const w = await api.updateWorkout(id, {
        status: 'completed',
        rpe,
        athleteReport,
      })
      setWorkout(w)
      setRpe(w.rpe ?? rpe)
      setAthleteReport(w.athleteReport || athleteReport)
      await reloadActivities()
      showOk('Тренировка отмечена выполненной')
    } catch {
      showErr('Не удалось отметить тренировку')
    } finally {
      setSaving(false)
    }
  }

  const saveReport = async () => {
    setSaving(true)
    setFlash(null)
    try {
      const w = await api.updateWorkout(id, { rpe, athleteReport })
      setWorkout(w)
      setRpe(w.rpe ?? rpe)
      setAthleteReport(w.athleteReport || athleteReport)
      showOk('Отчёт сохранён')
    } catch {
      showErr('Не удалось сохранить отчёт')
    } finally {
      setSaving(false)
    }
  }

  const saveCoachComment = async () => {
    setSaving(true)
    setFlash(null)
    try {
      const w = await api.updateWorkout(id, { coachComment })
      setWorkout(w)
      setCoachComment(w.coachComment || coachComment)
      showOk('Комментарий сохранён')
    } catch {
      showErr('Не удалось сохранить комментарий')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <button className="btn" onClick={closeOverlay} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>

      {flash && (
        <div
          data-testid="workout-flash"
          className="card"
          style={{
            background: flash.kind === 'ok' ? theme.accentSoft : theme.card2,
            color: flash.kind === 'ok' ? theme.accent : theme.accent,
            fontWeight: 700,
            fontSize: 13,
          }}
        >
          {flash.text}
        </div>
      )}

      <div className="card">
        <div style={{ fontSize: 11, color: theme.accent, fontWeight: 700 }}>{workoutTypeLabel(workout.workoutType)} · {workout.dayLabel}</div>
        <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 22, marginTop: 6 }}>{workout.title}</div>
        {workout.description && <div style={{ fontSize: 13, color: theme.dim, marginTop: 8 }}>{workout.description}</div>}
        <div style={{ display: 'flex', gap: 16, marginTop: 12, flexWrap: 'wrap' }}>
          <div><div style={{ fontWeight: 800 }}>{workout.distKm} км</div><div style={{ fontSize: 10, color: theme.dim }}>план</div></div>
          {workout.actualKm != null && (
            <div><div style={{ fontWeight: 800 }}>{workout.actualKm} км</div><div style={{ fontSize: 10, color: theme.dim }}>факт</div></div>
          )}
          <div><div style={{ fontWeight: 800 }}>{workout.status === 'completed' ? '✓' : '○'}</div><div style={{ fontSize: 10, color: theme.dim }}>статус</div></div>
          <div>
            <div style={{ fontWeight: 800 }}>{workout.rpe != null ? workout.rpe : '—'}</div>
            <div style={{ fontSize: 10, color: theme.dim }}>сложность</div>
          </div>
        </div>
      </div>

      {workout.segments?.map((s, i) => (
        <div key={s.id || i} className="card" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <div><div style={{ fontSize: 11, color: theme.accent }}>{segmentKindLabel(s.kind)}</div><div style={{ fontWeight: 700 }}>{s.title}</div></div>
          <div style={{ textAlign: 'right' }}><div style={{ fontWeight: 800 }}>{s.distKm} км</div><div style={{ fontSize: 11, color: theme.dim }}>{s.pace}</div></div>
        </div>
      ))}

      {workout.fact ? (
        <WorkoutFactPanel
          theme={theme}
          activity={workout.fact}
          onUpdated={(fact) => {
            setWorkout({ ...workout, fact })
            void reloadActivities()
            showOk('Данные тренировки обновлены')
          }}
        />
      ) : (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Ещё нет данных по активности</div>
      )}

      {canEditReport && (
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 14 }}>Отчёт</div>
          <div>
            <div style={{ fontSize: 11, color: theme.dim, marginBottom: 6 }}>Сложность (RPE) · {rpe}</div>
            <input
              type="range"
              min={1}
              max={10}
              value={rpe}
              onChange={(e) => setRpe(Number(e.target.value))}
              style={{ width: '100%', accentColor: theme.accent }}
            />
          </div>
          <textarea
            value={athleteReport}
            onChange={(e) => setAthleteReport(e.target.value)}
            placeholder="Как прошла тренировка…"
            rows={3}
            style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, resize: 'vertical', boxSizing: 'border-box' }}
          />
          {workout.status === 'completed' && (
            <button
              className="btn"
              disabled={saving}
              onClick={() => void saveReport()}
              style={{ background: theme.card2, color: theme.accent, borderRadius: 14, padding: 12 }}
            >
              Сохранить отчёт
            </button>
          )}
        </div>
      )}

      {isCoach && (
        <div className="card" data-testid="coach-difficulty" style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 14 }}>Отчёт ученика</div>
          <div>
            <div style={{ fontSize: 11, color: theme.dim, marginBottom: 6 }}>
              Сложность (RPE) · {workout.rpe != null ? workout.rpe : '—'}
            </div>
            <input
              type="range"
              min={1}
              max={10}
              value={workout.rpe ?? 1}
              disabled
              readOnly
              aria-label="Сложность тренировки"
              style={{
                width: '100%',
                accentColor: theme.accent,
                opacity: workout.rpe != null ? 1 : 0.35,
                cursor: 'default',
              }}
            />
            {workout.rpe == null && (
              <div style={{ fontSize: 12, color: theme.dim, marginTop: 4 }}>Ученик ещё не оценил сложность</div>
            )}
          </div>
          {workout.athleteReport ? (
            <div style={{ fontSize: 14 }}>{workout.athleteReport}</div>
          ) : (
            <div style={{ fontSize: 14, color: theme.dim }}>Пока нет текстового отчёта</div>
          )}
        </div>
      )}

      {isCoach ? (
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 14 }}>Комментарий тренера</div>
          <textarea
            value={coachComment}
            onChange={(e) => setCoachComment(e.target.value)}
            placeholder="Комментарий…"
            rows={3}
            style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, resize: 'vertical', boxSizing: 'border-box' }}
          />
          <button
            className="btn"
            disabled={saving}
            onClick={() => void saveCoachComment()}
            style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: saving ? 0.6 : 1 }}
          >
            {saving ? 'Сохранение…' : 'Сохранить комментарий'}
          </button>
        </div>
      ) : workout.coachComment ? (
        <div className="card">
          <div style={{ fontSize: 11, color: theme.dim, marginBottom: 6 }}>Комментарий тренера</div>
          <div style={{ fontSize: 14 }}>{workout.coachComment}</div>
        </div>
      ) : null}

      {canComplete && (
        <button
          className="btn"
          disabled={saving}
          onClick={() => void complete()}
          style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: saving ? 0.6 : 1 }}
        >
          Отметить выполненной
        </button>
      )}
      {workout.kind === 'own' && !isCoach && (
        <button className="btn" onClick={() => void api.deleteWorkout(id).then(closeOverlay)} style={{ background: theme.card2, color: theme.accent, borderRadius: 14, padding: 12 }}>
          Удалить
        </button>
      )}
    </div>
  )
}
