import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { SEGMENT_KINDS, segmentKindLabel, type SegmentKind } from '../../lib/workoutTypes'
import type { Theme } from '../../lib/theme'

export function SegmentKindField({ theme, value, onChange }: {
  theme: Theme
  value: string
  onChange: (kind: SegmentKind) => void
}) {
  const [open, setOpen] = useState(false)
  const [phoneEl, setPhoneEl] = useState<HTMLElement | null>(null)

  useEffect(() => {
    setPhoneEl(document.querySelector('.phone'))
  }, [])

  const sheet = open && phoneEl ? createPortal(
    <div
      data-testid="segment-kind-sheet"
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
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 18 }}>Тип сегмента</div>
          <button
            type="button"
            className="btn"
            onClick={() => setOpen(false)}
            style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}
          >
            ✕
          </button>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {SEGMENT_KINDS.map((k) => {
            const active = k.id === value
            return (
              <button
                key={k.id}
                type="button"
                data-testid="segment-kind-option"
                className="btn"
                onClick={() => {
                  onChange(k.id)
                  setOpen(false)
                }}
                style={{
                  width: '100%',
                  textAlign: 'left',
                  background: active ? theme.accentSoft : theme.card2,
                  border: `1.5px solid ${active ? theme.accent : 'transparent'}`,
                  borderRadius: 14,
                  padding: '12px 14px',
                  color: active ? theme.accent : theme.text,
                  fontWeight: 700,
                }}
              >
                {k.label}
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
        data-testid="segment-kind-field"
        className="btn"
        onClick={() => setOpen(true)}
        style={{
          background: theme.card2,
          color: theme.text,
          border: `1.5px solid ${open ? theme.accent : 'transparent'}`,
          borderRadius: 8,
          padding: '6px 10px',
          fontSize: 12,
          fontWeight: 700,
          display: 'flex',
          alignItems: 'center',
          gap: 6,
        }}
      >
        <span>{segmentKindLabel(value)}</span>
        <span style={{ color: theme.dim, fontSize: 10 }}>▾</span>
      </button>
      {sheet}
    </>
  )
}
