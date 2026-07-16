import { useEffect, useState } from 'react'
import { api, type StudentDetail } from '../lib/api'
import { useApp } from '../lib/store'
import { workoutTypeLabel } from '../lib/workoutTypes'

export function StudentProfileScreen({ id }: { id: string }) {
  const { theme, closeOverlay, openOverlay } = useApp()
  const [detail, setDetail] = useState<StudentDetail | null>(null)
  const [week, setWeek] = useState(0)

  useEffect(() => {
    void api.student(id, week).then(setDetail).catch(() => setDetail(null))
  }, [id, week])

  if (!detail) {
    return <div className="card" style={{ color: theme.dim }}>Загрузка…</div>
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      <button className="btn" onClick={closeOverlay} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>
      <div className="card" style={{ display: 'flex', gap: 14, alignItems: 'center' }}>
        <div style={{ width: 52, height: 52, borderRadius: '50%', background: theme.card2, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800, fontSize: 18 }}>{detail.student.init}</div>
        <div>
          <div style={{ fontWeight: 800, fontSize: 18 }}>{detail.student.name}</div>
          <div style={{ fontSize: 12, color: theme.dim }}>{detail.weekKm} / {detail.weekPlanKm} км · {detail.comp}%</div>
        </div>
      </div>
      <div style={{ display: 'flex', gap: 8 }}>
        {[0, 1, 2, 3].map((w) => (
          <button key={w} className="btn" onClick={() => setWeek(w)} style={{ flex: 1, padding: 8, borderRadius: 10, fontSize: 11, background: week === w ? theme.accent : theme.card2, color: week === w ? theme.onAccent : theme.dim }}>Нед {w + 1}</button>
        ))}
      </div>
      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase', fontSize: 14 }}>План vs факт</div>
      {detail.planDays.length === 0 && <div className="card" style={{ color: theme.dim, fontSize: 13 }}>Нет запланированных тренировок</div>}
      {detail.planDays.map((d, i) => (
        <div key={i} className="card" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <div>
            <div style={{ fontWeight: 700 }}>{d.dayLabel} · {d.title}</div>
            <div style={{ fontSize: 11, color: theme.dim }}>{workoutTypeLabel(d.workoutType)}</div>
          </div>
          <div style={{ textAlign: 'right' }}>
            <div style={{ fontWeight: 800, color: theme.accent }}>{d.actualKm}/{d.plannedKm} км</div>
            <div style={{ fontSize: 11, color: d.status === 'completed' ? theme.good : theme.dim }}>{d.status === 'completed' ? '✓' : '○'}</div>
          </div>
        </div>
      ))}
      <div style={{ fontFamily: theme.display, fontWeight: 800, textTransform: 'uppercase', fontSize: 14 }}>Активности</div>
      {detail.recentActivities.map((a) => (
        <div key={a.id} className="card" style={{ cursor: 'pointer' }} onClick={() => openOverlay({ type: 'activity', id: a.id })}>
          <div style={{ fontWeight: 700 }}>{a.title}</div>
          <div style={{ fontSize: 11, color: theme.dim }}>{a.when} · {a.dist} км · {a.pace}</div>
        </div>
      ))}
    </div>
  )
}
