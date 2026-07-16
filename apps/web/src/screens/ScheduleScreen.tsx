import { useEffect, useState } from 'react'
import { ScheduleCalendar, formatAnnounceDay } from '../components/schedule/ScheduleCalendar'
import { api } from '../lib/api'
import { useApp } from '../lib/store'

type CalCell = {
  key: string
  n: number
  blank: boolean
  has?: boolean
  isToday?: boolean
  iso?: string
  bg: string
  fg: string
  dot: string
}

function todayIso() {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

export function ScheduleScreen() {
  const { theme, announces, signupAnnounce, unsignupAnnounce, user, publishAnnounce, reloadAnnounces, setScreen } = useApp()
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [label, setLabel] = useState('')
  const [cells, setCells] = useState<CalCell[]>([])
  const [selectedIso, setSelectedIso] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [place, setPlace] = useState('')
  const [day, setDay] = useState('')
  const [time, setTime] = useState('19:50')
  const [group, setGroup] = useState('Основная группа')
  const [note, setNote] = useState('')
  const [saving, setSaving] = useState(false)

  const loadCalendar = (y: number, m: number) => {
    void api.calendar(y, m).then((r) => {
      setCells(r.cells)
      setLabel(r.label)
      setYear(r.year || y)
      setMonth(r.month || m)
    }).catch(() => {
      setCells([])
      setLabel('')
    })
  }

  useEffect(() => {
    void reloadAnnounces()
  }, [reloadAnnounces])

  useEffect(() => {
    loadCalendar(year, month)
  }, [year, month])

  useEffect(() => {
    if (selectedIso) setDay(formatAnnounceDay(selectedIso))
  }, [selectedIso])

  const shiftMonth = (delta: number) => {
    const d = new Date(year, month - 1 + delta, 1)
    setYear(d.getFullYear())
    setMonth(d.getMonth() + 1)
  }

  const selectDay = (iso: string, dayLabel: string) => {
    setSelectedIso(iso)
    setDay(dayLabel)
    if (user?.role === 'coach') setShowForm(true)
  }

  const ensureDateSelected = () => {
    if (selectedIso) return selectedIso
    const fromCells = cells.find((c) => c.isToday && c.iso)?.iso
    const iso = fromCells || todayIso()
    setSelectedIso(iso)
    setDay(formatAnnounceDay(iso))
    return iso
  }

  const toggleForm = () => {
    if (showForm) {
      setShowForm(false)
      return
    }
    ensureDateSelected()
    setShowForm(true)
  }

  const publish = async () => {
    const iso = selectedIso || ensureDateSelected()
    if (!place.trim() || !iso) return
    const dayLabel = day.trim() || formatAnnounceDay(iso)
    setSaving(true)
    try {
      await publishAnnounce({
        place: place.trim(),
        day: dayLabel,
        time: time.trim() || '19:50',
        group: group.trim() || 'Основная группа',
        note: note.trim(),
        startsOn: iso,
      })
      setShowForm(false)
      setPlace('')
      setDay('')
      setNote('')
      setSelectedIso(null)
      loadCalendar(year, month)
    } finally {
      setSaving(false)
    }
  }

  const visibleAnnounces = selectedIso
    ? announces.filter((a) => a.startsOn === selectedIso)
    : announces

  const canPublish = Boolean(place.trim() && selectedIso)

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
      <button className="btn" onClick={() => setScreen('home')} style={{ alignSelf: 'flex-start', background: theme.card2, color: theme.dim, borderRadius: 10, padding: '8px 12px' }}>← Назад</button>
      <ScheduleCalendar
        theme={theme}
        label={label}
        cells={cells}
        selectedIso={selectedIso}
        onSelect={selectDay}
        onPrevMonth={() => shiftMonth(-1)}
        onNextMonth={() => shiftMonth(1)}
      />

      {user?.role === 'coach' && (
        <button className="btn" onClick={toggleForm} style={{ border: `1px solid ${theme.accent}`, background: theme.accentSoft, color: theme.accent, borderRadius: 14, padding: 13 }}>
          ＋ Опубликовать анонс
        </button>
      )}
      {showForm && (
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <input value={place} onChange={(e) => setPlace(e.target.value)} placeholder="Место" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }} />
          <div style={{ width: '100%', background: theme.card2, borderRadius: 10, padding: 12, color: selectedIso ? theme.text : theme.dim, boxSizing: 'border-box', fontSize: 14 }}>
            {selectedIso ? formatAnnounceDay(selectedIso) : 'Выберите день в календаре'}
          </div>
          <input value={time} onChange={(e) => setTime(e.target.value)} placeholder="Время" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }} />
          <input value={group} onChange={(e) => setGroup(e.target.value)} placeholder="Группа" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }} />
          <input value={note} onChange={(e) => setNote(e.target.value)} placeholder="Заметка" style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }} />
          <button
            data-testid="publish-announce"
            className="btn"
            disabled={saving || !canPublish}
            onClick={() => void publish()}
            style={{ background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: canPublish ? 1 : 0.5 }}
          >
            Опубликовать для учеников
          </button>
        </div>
      )}

      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 8 }}>
        <span style={{ fontFamily: theme.display, fontSize: 16, fontWeight: 800, textTransform: 'uppercase' }}>
          {selectedIso ? formatAnnounceDay(selectedIso) : 'Ближайшие занятия'}
        </span>
        {selectedIso ? (
          <button
            type="button"
            className="btn"
            onClick={() => setSelectedIso(null)}
            style={{ background: theme.card2, color: theme.dim, borderRadius: 10, padding: '6px 10px', fontSize: 12 }}
          >
            Сбросить
          </button>
        ) : null}
      </div>
      {visibleAnnounces.length === 0 && (
        <div className="card" style={{ fontSize: 13, color: theme.dim }}>
          {selectedIso ? 'Нет анонсов на этот день' : 'Пока нет анонсов'}
        </div>
      )}
      {visibleAnnounces.map((an) => (
        <div key={an.id} className="card" data-testid="announce-card" style={{ position: 'relative', overflow: 'hidden' }}>
          <div style={{ position: 'absolute', left: 0, top: 0, bottom: 0, width: 4, background: theme.accent }} />
          <div style={{ fontSize: 10, fontWeight: 800, color: theme.accent }}>ГРУППОВАЯ · {an.group}</div>
          <div style={{ fontFamily: theme.display, fontSize: 19, fontWeight: 800, margin: '6px 0' }}>{an.place}</div>
          <div style={{ fontSize: 13, color: theme.dim, marginBottom: 10 }}>{an.day} · {an.time} · идут {an.going}</div>
          {(an.attendees?.length ?? 0) > 0 && (
            <div data-testid="announce-attendees" style={{ display: 'flex', flexDirection: 'column', gap: 8, marginBottom: 14 }}>
              <div style={{ display: 'flex', alignItems: 'center' }}>
                {an.attendees!.slice(0, 6).map((p, i) => (
                  <div
                    key={p.id}
                    title={p.name}
                    style={{
                      width: 32,
                      height: 32,
                      borderRadius: '50%',
                      background: theme.card2,
                      border: `2px solid ${theme.card}`,
                      marginLeft: i === 0 ? 0 : -8,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: theme.accent,
                      fontSize: 10,
                      fontWeight: 800,
                      zIndex: 6 - i,
                    }}
                  >
                    {p.init}
                  </div>
                ))}
                {(an.attendees!.length > 6) && (
                  <span style={{ marginLeft: 8, fontSize: 11, color: theme.dim, fontWeight: 700 }}>
                    +{an.attendees!.length - 6}
                  </span>
                )}
              </div>
              <div style={{ fontSize: 12, color: theme.dim }}>
                {an.attendees!.map((p) => p.name).join(', ')}
              </div>
            </div>
          )}
          {user?.role === 'athlete' && (
            <button
              data-testid="schedule-cta"
              className="btn"
              onClick={() => void (an.signedUp ? unsignupAnnounce(an.id) : signupAnnounce(an.id))}
              style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 11, padding: 12, fontSize: 13 }}
            >
              {an.scheduleCta}
            </button>
          )}
        </div>
      ))}
    </div>
  )
}
