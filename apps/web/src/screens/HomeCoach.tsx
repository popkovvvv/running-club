import { useEffect, useState } from 'react'
import { api, type Student } from '../lib/api'
import { useApp } from '../lib/store'

export function HomeCoach() {
  const { theme, announces, removeStudent, setScreen } = useApp()
  const [students, setStudents] = useState<Student[]>([])

  useEffect(() => {
    void api.students().then(setStudents).catch(() => setStudents([]))
  }, [])

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <div style={{ display: 'flex', gap: 10 }}>
        <div className="card" style={{ flex: 1, borderRadius: 16 }}><div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800 }}>{students.length}</div><div style={{ fontSize: 11, color: theme.dim }}>учеников</div></div>
        <div className="card" style={{ flex: 1, borderRadius: 16 }}><div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800, color: theme.accent }}>3</div><div style={{ fontSize: 11, color: theme.dim }}>сегодня</div></div>
        <div className="card" style={{ flex: 1, borderRadius: 16 }}><div style={{ fontFamily: theme.display, fontSize: 26, fontWeight: 800 }}>86%</div><div style={{ fontSize: 11, color: theme.dim }}>явка</div></div>
      </div>
      <div style={{ fontFamily: theme.display, fontSize: 19, fontWeight: 800, textTransform: 'uppercase' }}>Анонсы клуба</div>
      {announces.map((an) => (
        <div key={an.id} className="card" style={{ borderRadius: 16, display: 'flex', gap: 13, alignItems: 'center' }}>
          <div style={{ width: 4, height: 40, borderRadius: 3, background: theme.accent }} />
          <div style={{ flex: 1 }}><div style={{ fontWeight: 700 }}>{an.place}</div><div style={{ fontSize: 11, color: theme.dim }}>{an.day} · {an.time}</div></div>
          <div style={{ textAlign: 'right' }}><div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{an.going}</div><div style={{ fontSize: 10, color: theme.dim }}>идут</div></div>
        </div>
      ))}
      <div style={{ fontFamily: theme.display, fontSize: 19, fontWeight: 800, textTransform: 'uppercase' }}>Мои ученики</div>
      {students.map((s) => (
        <div key={s.id} className="card" data-testid="student-row" style={{ borderRadius: 16, display: 'flex', gap: 13, alignItems: 'center' }}>
          <div style={{ width: 46, height: 46, borderRadius: '50%', background: theme.card2, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800 }}>{s.init}</div>
          <div style={{ flex: 1 }} onClick={() => setScreen('plan')}><div style={{ fontWeight: 700 }}>{s.name}</div><div style={{ fontSize: 11, color: theme.dim }}>{s.sub}</div></div>
          <button
            data-testid="remove-student"
            className="btn"
            onClick={() => void removeStudent(s.id).then(() => setStudents((x) => x.filter((i) => i.id !== s.id)))}
            style={{ background: theme.card2, color: theme.accent, borderRadius: 10, padding: '8px 10px', fontSize: 12 }}
          >
            Удалить
          </button>
        </div>
      ))}
    </div>
  )
}
