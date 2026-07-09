import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { useApp } from '../lib/store'

export function PlanScreen() {
  const { theme, user } = useApp()
  const [plan, setPlan] = useState<Awaited<ReturnType<typeof api.plan>> | null>(null)
  const [week, setWeek] = useState(1)
  const [open, setOpen] = useState(false)
  const [segs, setSegs] = useState([
    { kind: 'Разминка', title: 'Лёгкий бег', distKm: 2, pace: '7:40' },
    { kind: 'Основная', title: 'ОФП', distKm: 5, pace: 'смеш.' },
    { kind: 'Заминка', title: 'Растяжка', distKm: 1, pace: '8:00' },
  ])

  useEffect(() => {
    void api.plan(week).then(setPlan)
  }, [week])

  const total = segs.reduce((a, s) => a + s.distKm, 0)

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      {user?.role === 'athlete' ? (
        <>
          <div className="card" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', borderRadius: 16 }}>
            <button className="btn" onClick={() => setWeek((w) => Math.max(0, w - 1))} style={{ width: 38, height: 38, background: theme.card2, color: theme.text, borderRadius: 11 }}>‹</button>
            <div style={{ textAlign: 'center' }}><div style={{ fontWeight: 700 }}>{plan?.weekRange || '…'}</div><div style={{ fontSize: 11, color: theme.accent }}>План {plan?.weekPlan}</div></div>
            <button className="btn" onClick={() => setWeek((w) => Math.min(3, w + 1))} style={{ width: 38, height: 38, background: theme.card2, color: theme.text, borderRadius: 11 }}>›</button>
          </div>
          <button data-testid="create-workout" className="btn" onClick={() => setOpen(true)} style={{ border: `1.5px dashed ${theme.accent}`, background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13 }}>
            ＋ Создать свою тренировку
          </button>
          {plan?.mine?.map((w) => (
            <div key={w.id} className="card" style={{ borderColor: theme.accent }}>
              <div style={{ fontSize: 11, color: theme.accent, fontWeight: 700 }}>{w.tag}</div>
              <div style={{ fontWeight: 700 }}>{w.title}</div>
              <div style={{ fontSize: 12, color: theme.dim }}>{w.dayLabel} · {w.duration}</div>
            </div>
          ))}
          {plan?.days?.map((d) => (
            <div key={d.id} className="card" style={{ borderRadius: 16 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <div><div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{d.dayLabel}</div><div style={{ fontSize: 11, color: theme.dim }}>{d.tag}</div></div>
                <div style={{ textAlign: 'right' }}><div style={{ fontWeight: 800 }}>{d.distKm || '—'} км</div><div style={{ fontSize: 11, color: theme.accent }}>{d.duration}</div></div>
              </div>
              <div style={{ fontWeight: 700, marginTop: 8 }}>{d.title}</div>
            </div>
          ))}
          {open && (
            <div style={{ position: 'absolute', inset: 0, background: 'rgba(0,0,0,.55)', zIndex: 20, display: 'flex', alignItems: 'flex-end' }}>
              <div className="card" style={{ width: '100%', borderRadius: '28px 28px 46px 46px' }}>
                <div style={{ fontFamily: theme.display, fontWeight: 800, marginBottom: 12 }}>Новая тренировка</div>
                <button
                  data-testid="save-workout"
                  className="btn"
                  onClick={() => void api.createWorkout({ kind: 'own', dayLabel: 'Ср', tag: 'Лёгкий', title: 'Свой кросс', distKm: 6, duration: '~46 мин', pace: '7:40', weekIndex: week }).then(() => api.plan(week).then(setPlan)).then(() => setOpen(false))}
                  style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14 }}
                >
                  Добавить в план
                </button>
                <button className="btn" onClick={() => setOpen(false)} style={{ width: '100%', marginTop: 8, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}>Закрыть</button>
              </div>
            </div>
          )}
        </>
      ) : (
        <>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 19, textTransform: 'uppercase' }}>Конструктор</div>
          <div className="card">Объём: <b style={{ color: theme.accent }}>{total} км</b> · сегментов {segs.length}</div>
          {segs.map((s, i) => (
            <div key={i} className="card" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div><div style={{ fontSize: 11, color: theme.accent }}>{s.kind}</div><div style={{ fontWeight: 700 }}>{s.title}</div></div>
              <button className="btn" onClick={() => setSegs((arr) => arr.filter((_, j) => j !== i))} style={{ background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 10px' }}>✕</button>
            </div>
          ))}
          <button className="btn" onClick={() => setSegs((arr) => [...arr, { kind: 'Интервалы', title: '5×800', distKm: 4, pace: '4:30' }])} style={{ background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13, border: `1px dashed ${theme.accent}` }}>＋ Сегмент</button>
          <button
            className="btn"
            onClick={() => void api.createWorkout({ kind: 'builder', title: 'Групповая', segments: segs, weekIndex: week })}
            style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 15 }}
          >
            Назначить ученику
          </button>
        </>
      )}
    </div>
  )
}
