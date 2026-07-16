import { useEffect, useState } from 'react'
import { AccentPill } from '../components/ui/AccentPill'
import { HeroPanel } from '../components/ui/HeroPanel'
import { api, type Student } from '../lib/api'
import { nearestAnnounce } from '../lib/nearestAnnounce'
import { useApp } from '../lib/store'

export function HomeCoach() {
  const { theme, announces, removeStudent, openOverlay, setScreen, reloadAnnounces, tabEpoch } = useApp()
  const [students, setStudents] = useState<Student[]>([])

  useEffect(() => {
    void reloadAnnounces()
    void api.students().then(setStudents).catch(() => setStudents([]))
  }, [tabEpoch, reloadAnnounces])

  const next = nearestAnnounce(announces)

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <HeroPanel>
        <div className="display-title" style={{ fontSize: 26, marginBottom: 10 }}>Новая неделя</div>
        {next ? (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, alignItems: 'center', marginBottom: 8 }}>
            <AccentPill>{next.day}</AccentPill>
            <span style={{ fontFamily: theme.display, fontWeight: 800, fontStyle: 'italic' }}>{next.time}</span>
          </div>
        ) : (
          <div style={{ fontSize: 13, color: theme.dim, marginBottom: 8 }}>Пока нет анонсов</div>
        )}
        <div style={{ display: 'flex', gap: 16, marginTop: 6 }}>
          <div>
            <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800 }}>{students.length}</div>
            <div style={{ fontSize: 11, color: theme.dim }}>учеников</div>
          </div>
          <div>
            <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800, color: theme.accent }}>{announces.length}</div>
            <div style={{ fontSize: 11, color: theme.dim }}>анонсов</div>
          </div>
        </div>
      </HeroPanel>
      <button className="btn" onClick={() => setScreen('schedule')} style={{ border: `1.5px dashed ${theme.accent}`, background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13 }}>
        ＋ Опубликовать анонс / расписание
      </button>
      <div className="display-title" style={{ fontSize: 18 }}>Анонсы клуба</div>
      {announces.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Нет анонсов</div>
      )}
      {announces.map((an) => (
        <div key={an.id} className="card" style={{ borderRadius: 16, display: 'flex', flexDirection: 'column', gap: 10, cursor: 'pointer' }} onClick={() => setScreen('schedule')}>
          <div style={{ display: 'flex', gap: 13, alignItems: 'center' }}>
            <div style={{ width: 4, height: 40, borderRadius: 3, background: theme.accent }} />
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 700 }}>{an.place}</div>
              <div style={{ display: 'flex', gap: 8, alignItems: 'center', marginTop: 6 }}>
                <AccentPill style={{ fontSize: 10, padding: '4px 10px' }}>{an.day}</AccentPill>
                <span style={{ fontSize: 11, color: theme.dim }}>{an.time}</span>
              </div>
            </div>
            <div style={{ textAlign: 'right' }}>
              <div style={{ fontFamily: theme.display, fontWeight: 800, color: theme.accent }}>{an.going}</div>
              <div style={{ fontSize: 10, color: theme.dim }}>идут</div>
            </div>
          </div>
          {(an.attendees?.length ?? 0) > 0 && (
            <div style={{ display: 'flex', alignItems: 'center', paddingLeft: 17 }}>
              {an.attendees!.slice(0, 5).map((p, i) => (
                <div
                  key={p.id}
                  title={p.name}
                  style={{
                    width: 28,
                    height: 28,
                    borderRadius: '50%',
                    background: theme.card2,
                    border: `2px solid ${theme.card}`,
                    marginLeft: i === 0 ? 0 : -6,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    color: theme.accent,
                    fontSize: 9,
                    fontWeight: 800,
                  }}
                >
                  {p.init}
                </div>
              ))}
              <span style={{ marginLeft: 8, fontSize: 11, color: theme.dim }}>
                {an.attendees!.map((p) => p.name.split(' ')[0]).join(', ')}
              </span>
            </div>
          )}
        </div>
      ))}
      <div className="display-title" style={{ fontSize: 18 }}>Мои ученики</div>
      {students.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>Пока нет учеников</div>
      )}
      {students.map((s) => (
        <div key={s.id} className="card" data-testid="student-row" style={{ borderRadius: 16, display: 'flex', gap: 13, alignItems: 'center' }}>
          <div style={{ width: 46, height: 46, borderRadius: '50%', background: theme.card2, display: 'flex', alignItems: 'center', justifyContent: 'center', color: theme.accent, fontWeight: 800 }}>{s.init}</div>
          <div style={{ flex: 1, cursor: 'pointer' }} onClick={() => openOverlay({ type: 'student', id: s.id })}>
            <div style={{ fontWeight: 700 }}>{s.name}</div>
            <div style={{ fontSize: 11, color: theme.dim }}>{s.km} км · {s.comp}%</div>
          </div>
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
