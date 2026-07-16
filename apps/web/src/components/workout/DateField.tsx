import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import type { Theme } from '../../lib/theme'
import { DatePicker } from './DatePicker'

export function DateField({ theme, value, dayLabel, onChange }: {
  theme: Theme
  value: string | null
  dayLabel: string
  onChange: (iso: string, dayLabel: string) => void
}) {
  const [open, setOpen] = useState(false)
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  const label = value
    ? `${dayLabel}, ${value.split('-').reverse().join('.')}`
    : 'Выберите дату'

  const sheet = open && phoneEl ? createPortal(
    <div
      data-testid="date-sheet"
      style={{
        position: 'absolute',
        inset: 0,
        background: 'rgba(0,0,0,.55)',
        zIndex: 60,
        display: 'flex',
        alignItems: 'flex-end',
      }}
      onClick={() => setOpen(false)}
    >
      <div
        className="card"
        onClick={(e) => e.stopPropagation()}
        style={{
          width: '100%',
          borderRadius: '24px 24px 0 0',
          paddingBottom: 28,
          boxSizing: 'border-box',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 14 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 18 }}>Дата тренировки</div>
          <button
            type="button"
            className="btn"
            data-testid="date-sheet-close"
            onClick={() => setOpen(false)}
            style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}
          >
            ✕
          </button>
        </div>
        <DatePicker
          theme={theme}
          value={value}
          onChange={(iso, dLabel) => {
            onChange(iso, dLabel)
            setOpen(false)
          }}
        />
        <button
          type="button"
          className="btn"
          onClick={() => setOpen(false)}
          style={{ width: '100%', marginTop: 14, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}
        >
          Закрыть
        </button>
      </div>
    </div>,
    phoneEl,
  ) : null

  return (
    <>
      <button
        type="button"
        data-testid="date-field"
        className="btn"
        onClick={() => setOpen(true)}
        style={{
          width: '100%',
          textAlign: 'left',
          background: theme.card2,
          border: `1.5px solid ${open ? theme.accent : 'transparent'}`,
          borderRadius: 10,
          padding: 12,
          color: value ? theme.text : theme.dim,
          boxSizing: 'border-box',
          fontWeight: 600,
        }}
      >
        {label}
      </button>
      {sheet}
    </>
  )
}
