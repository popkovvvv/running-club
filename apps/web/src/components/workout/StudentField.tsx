import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import type { Student } from '../../lib/api'
import type { Theme } from '../../lib/theme'
import { AccentPill } from '../ui/AccentPill'

export function StudentField({ theme, students, value, onChange }: {
  theme: Theme
  students: Student[]
  value: string
  onChange: (id: string) => void
}) {
  const [open, setOpen] = useState(false)
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)
  const selected = students.find((s) => s.id === value)

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  const sheet = open && phoneEl ? createPortal(
    <div
      data-testid="student-sheet"
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
          maxHeight: '78%',
          overflowY: 'auto',
          borderRadius: '24px 24px 0 0',
          paddingBottom: 28,
          boxSizing: 'border-box',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 14 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 18 }}>Ученик</div>
          <button
            type="button"
            className="btn"
            data-testid="student-sheet-close"
            onClick={() => setOpen(false)}
            style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}
          >
            ✕
          </button>
        </div>
        {students.length === 0 && (
          <div style={{ fontSize: 13, color: theme.dim, padding: '8px 0' }}>Пока нет учеников</div>
        )}
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {students.map((s) => {
            const active = s.id === value
            return (
              <button
                key={s.id}
                type="button"
                data-testid="student-option"
                className="btn"
                onClick={() => {
                  onChange(s.id)
                  setOpen(false)
                }}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 12,
                  width: '100%',
                  textAlign: 'left',
                  background: active ? theme.accentSoft : theme.card2,
                  border: `1.5px solid ${active ? theme.accent : 'transparent'}`,
                  borderRadius: 14,
                  padding: '12px 14px',
                  color: theme.text,
                }}
              >
                <div style={{
                  width: 40,
                  height: 40,
                  borderRadius: '50%',
                  background: active ? theme.accent : theme.bg,
                  color: active ? theme.onAccent : theme.accent,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontWeight: 800,
                  fontSize: 13,
                  flex: 'none',
                }}
                >
                  {s.init}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ fontWeight: 700 }}>{s.name}</div>
                  <div style={{ fontSize: 11, color: theme.dim, marginTop: 2 }}>{s.km} км · {s.comp}%</div>
                </div>
                {active ? <AccentPill style={{ fontSize: 10, padding: '4px 10px' }}>Выбран</AccentPill> : null}
              </button>
            )
          })}
        </div>
      </div>
    </div>,
    phoneEl,
  ) : null

  return (
    <>
      <button
        type="button"
        data-testid="student-field"
        className="btn"
        onClick={() => setOpen(true)}
        style={{
          width: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 10,
          textAlign: 'left',
          background: theme.card2,
          border: `1.5px solid ${open ? theme.accent : 'transparent'}`,
          borderRadius: 10,
          padding: 12,
          color: selected ? theme.text : theme.dim,
          boxSizing: 'border-box',
          fontWeight: 600,
        }}
      >
        <span style={{ display: 'flex', alignItems: 'center', gap: 10, minWidth: 0 }}>
          {selected ? (
            <span style={{
              width: 28,
              height: 28,
              borderRadius: '50%',
              background: theme.accent,
              color: theme.onAccent,
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 11,
              fontWeight: 800,
              flex: 'none',
            }}
            >
              {selected.init}
            </span>
          ) : null}
          <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
            {selected ? selected.name : 'Выберите ученика'}
          </span>
        </span>
        <span style={{ color: theme.dim, fontSize: 12 }}>▾</span>
      </button>
      {sheet}
    </>
  )
}
