import { useState } from 'react'
import type { Theme } from '../../lib/theme'

const WEEKDAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
const MONTHS = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь']

function startOfMonth(d: Date) {
  return new Date(d.getFullYear(), d.getMonth(), 1)
}

function daysInMonth(d: Date) {
  return new Date(d.getFullYear(), d.getMonth() + 1, 0).getDate()
}

function mondayIndex(d: Date) {
  const day = d.getDay()
  return day === 0 ? 6 : day - 1
}

function fmt(d: Date) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function dayLabelRu(d: Date) {
  return WEEKDAYS[mondayIndex(d)]
}

export function DatePicker({ theme, value, onChange }: {
  theme: Theme
  value: string | null
  onChange: (iso: string, dayLabel: string) => void
}) {
  const selected = value ? new Date(value + 'T12:00:00') : null
  const [cursor, setCursor] = useState(selected || new Date())

  const first = startOfMonth(cursor)
  const blanks = mondayIndex(first)
  const total = daysInMonth(cursor)
  const cells: Array<{ n: number; blank: boolean; iso?: string }> = []
  for (let i = 0; i < blanks; i++) cells.push({ n: 0, blank: true })
  for (let n = 1; n <= total; n++) {
    const d = new Date(cursor.getFullYear(), cursor.getMonth(), n)
    cells.push({ n, blank: false, iso: fmt(d) })
  }

  return (
    <div style={{ background: theme.card2, borderRadius: 14, padding: 12 }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 10 }}>
        <button
          type="button"
          className="btn"
          onClick={() => setCursor(new Date(cursor.getFullYear(), cursor.getMonth() - 1, 1))}
          style={{ width: 32, height: 32, background: theme.card, color: theme.text, borderRadius: 8 }}
        >
          ‹
        </button>
        <div style={{ fontWeight: 800, fontSize: 13 }}>{MONTHS[cursor.getMonth()]} {cursor.getFullYear()}</div>
        <button
          type="button"
          className="btn"
          onClick={() => setCursor(new Date(cursor.getFullYear(), cursor.getMonth() + 1, 1))}
          style={{ width: 32, height: 32, background: theme.card, color: theme.text, borderRadius: 8 }}
        >
          ›
        </button>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 4, marginBottom: 4 }}>
        {WEEKDAYS.map((w) => (
          <div key={w} style={{ textAlign: 'center', fontSize: 10, color: theme.dim, fontWeight: 700 }}>{w}</div>
        ))}
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: 4 }}>
        {cells.map((c, i) => {
          const active = !c.blank && c.iso === value
          return (
            <button
              key={i}
              type="button"
              className="btn"
              disabled={c.blank}
              onClick={() => {
                if (!c.iso) return
                const d = new Date(c.iso + 'T12:00:00')
                onChange(c.iso, dayLabelRu(d))
              }}
              style={{
                aspectRatio: '1',
                borderRadius: active ? 999 : 9,
                border: 'none',
                background: active ? theme.accent : c.blank ? 'transparent' : theme.card,
                color: active ? theme.onAccent : theme.text,
                fontSize: 12,
                fontWeight: 700,
                fontStyle: active ? 'italic' : 'normal',
                opacity: c.blank ? 0 : 1,
                cursor: c.blank ? 'default' : 'pointer',
              }}
            >
              {c.n || ''}
            </button>
          )
        })}
      </div>
    </div>
  )
}
