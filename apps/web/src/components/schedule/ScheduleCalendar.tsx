import type { Theme } from '../../lib/theme'

const WEEKDAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
const MONTHS_GENITIVE = [
  'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря',
]

export type ScheduleCalCell = {
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

export function formatAnnounceDay(iso: string): string {
  const d = new Date(iso + 'T12:00:00')
  const wd = WEEKDAYS[(d.getDay() + 6) % 7]
  return `${wd}, ${d.getDate()} ${MONTHS_GENITIVE[d.getMonth()]}`
}

export function ScheduleCalendar({
  theme,
  label,
  cells,
  selectedIso,
  onSelect,
  onPrevMonth,
  onNextMonth,
}: {
  theme: Theme
  label: string
  cells: ScheduleCalCell[]
  selectedIso: string | null
  onSelect: (iso: string, dayLabel: string) => void
  onPrevMonth: () => void
  onNextMonth: () => void
}) {
  return (
    <div className="card" data-testid="schedule-calendar">
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
        <button
          type="button"
          className="btn"
          data-testid="cal-prev"
          onClick={onPrevMonth}
          style={{ width: 36, height: 36, background: theme.card2, color: theme.text, borderRadius: 10 }}
        >
          ‹
        </button>
        <div style={{ fontFamily: theme.display, fontWeight: 800 }}>{label || '…'}</div>
        <button
          type="button"
          className="btn"
          data-testid="cal-next"
          onClick={onNextMonth}
          style={{ width: 36, height: 36, background: theme.card2, color: theme.text, borderRadius: 10 }}
        >
          ›
        </button>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 6, marginBottom: 6 }}>
        {WEEKDAYS.map((w) => (
          <div key={w} style={{ textAlign: 'center', fontSize: 10, color: theme.dim, fontWeight: 700 }}>{w}</div>
        ))}
      </div>
      <div data-testid="calendar-grid" style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 6 }}>
        {cells.map((c) => {
          const selected = !c.blank && !!c.iso && c.iso === selectedIso
          let bg = c.bg
          let fg = c.fg
          if (selected) {
            bg = theme.accent
            fg = theme.onAccent
          } else if (c.isToday) {
            bg = theme.accentSoft
            fg = theme.text
          }
          const showDot = !c.blank && !!c.has && !selected
          return (
            <button
              key={c.key}
              type="button"
              className="btn"
              data-testid={c.blank ? 'cal-blank' : 'cal-day'}
              data-blank={c.blank ? '1' : '0'}
              data-iso={c.iso || undefined}
              disabled={c.blank}
              onClick={() => {
                if (!c.iso) return
                onSelect(c.iso, formatAnnounceDay(c.iso))
              }}
              style={{
                aspectRatio: '1',
                borderRadius: selected ? 999 : 9,
                background: bg,
                color: fg,
                border: !selected && c.isToday ? `1.5px solid ${theme.accent}` : 'none',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative',
                opacity: c.blank ? 0 : 1,
                cursor: c.blank ? 'default' : 'pointer',
                fontStyle: selected ? 'italic' : 'normal',
                padding: 0,
              }}
            >
              <span style={{ fontSize: 12, fontWeight: 700 }}>{c.n || ''}</span>
              {showDot ? (
                <span
                  data-testid="cal-dot"
                  style={{ width: 5, height: 5, borderRadius: '50%', background: theme.accent, position: 'absolute', bottom: 5 }}
                />
              ) : null}
            </button>
          )
        })}
      </div>
    </div>
  )
}
