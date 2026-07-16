import { useState } from 'react'
import { createPortal } from 'react-dom'
import type { ActivityDetail } from '../../lib/api'
import type { Theme } from '../../lib/theme'
import { DateField } from '../workout/DateField'

export type ActivityEditValues = {
  title: string
  when: string
  dayLabel: string
  distKm: string
  duration: string
  pace: string
  hr: string
  elevationGain: string
}

const WEEKDAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

function dayLabelFromIso(iso: string): string {
  const d = new Date(iso + 'T12:00:00')
  if (Number.isNaN(d.getTime())) return ''
  const day = d.getDay()
  const idx = day === 0 ? 6 : day - 1
  return WEEKDAYS[idx]
}

function toIsoDate(activity: ActivityDetail): string {
  if (activity.startedAt) {
    const d = new Date(activity.startedAt)
    if (!Number.isNaN(d.getTime())) {
      const y = d.getFullYear()
      const m = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      return `${y}-${m}-${day}`
    }
  }
  const m = activity.when.match(/(\d{2})\.(\d{2})\.(\d{4})/)
  if (m) return `${m[3]}-${m[2]}-${m[1]}`
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
}

function fromActivity(activity: ActivityDetail): ActivityEditValues {
  const when = toIsoDate(activity)
  return {
    title: activity.title || '',
    when,
    dayLabel: dayLabelFromIso(when),
    distKm: activity.dist ? String(parseFloat(activity.dist) || '') : '',
    duration: activity.time || '',
    pace: activity.pace || '',
    hr: activity.hr && activity.hr !== '0' ? activity.hr : '',
    elevationGain: activity.elevation ? String(Math.round(activity.elevation)) : '',
  }
}

export function ActivityEditForm({
  theme,
  activity,
  phoneEl,
  onClose,
  onSave,
}: {
  theme: Theme
  activity: ActivityDetail
  phoneEl: HTMLElement
  onClose: () => void
  onSave: (body: Record<string, unknown>) => Promise<void>
}) {
  const [values, setValues] = useState(() => fromActivity(activity))
  const [saving, setSaving] = useState(false)

  const set = (key: keyof ActivityEditValues, value: string) => {
    setValues((v) => ({ ...v, [key]: value }))
  }

  const submit = async () => {
    setSaving(true)
    try {
      const body: Record<string, unknown> = {
        title: values.title.trim(),
        when: values.when,
        duration: values.duration.trim(),
        pace: values.pace.trim(),
      }
      const dist = parseFloat(values.distKm.replace(',', '.'))
      if (!Number.isNaN(dist)) body.distKm = dist
      const hr = parseInt(values.hr, 10)
      if (!Number.isNaN(hr)) body.hr = hr
      const elev = parseFloat(values.elevationGain.replace(',', '.'))
      if (!Number.isNaN(elev)) body.elevationGain = elev
      await onSave(body)
      onClose()
    } finally {
      setSaving(false)
    }
  }

  const field = (label: string, key: Exclude<keyof ActivityEditValues, 'when' | 'dayLabel'>, opts?: { placeholder?: string }) => (
    <label style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
      <span style={{ fontSize: 12, color: theme.dim, fontWeight: 700 }}>{label}</span>
      <input
        type="text"
        value={values[key]}
        placeholder={opts?.placeholder}
        onChange={(e) => set(key, e.target.value)}
        style={{ width: '100%', background: theme.card2, border: 'none', borderRadius: 10, padding: 12, color: theme.text, boxSizing: 'border-box' }}
      />
    </label>
  )

  return createPortal(
    <div
      data-testid="activity-edit-modal"
      style={{
        position: 'absolute',
        inset: 0,
        background: 'rgba(0,0,0,.6)',
        zIndex: 50,
        display: 'flex',
        alignItems: 'flex-end',
      }}
      onClick={onClose}
    >
      <div
        className="card"
        onClick={(e) => e.stopPropagation()}
        style={{
          width: '100%',
          maxHeight: '92%',
          overflowY: 'auto',
          borderRadius: '24px 24px 0 0',
          paddingBottom: 28,
          boxSizing: 'border-box',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800 }}>Редактировать</div>
          <button type="button" className="btn" onClick={onClose} style={{ width: 36, height: 36, borderRadius: 10, background: theme.card2, color: theme.dim, fontSize: 18 }}>✕</button>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
          {field('Название', 'title', { placeholder: 'Название' })}
          <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
            <span style={{ fontSize: 12, color: theme.dim, fontWeight: 700 }}>Дата</span>
            <DateField
              theme={theme}
              value={values.when}
              dayLabel={values.dayLabel}
              onChange={(iso, dayLabel) => setValues((v) => ({ ...v, when: iso, dayLabel }))}
            />
          </div>
          {field('Дистанция, км', 'distKm', { placeholder: '10.5' })}
          {field('Время', 'duration', { placeholder: '46:30' })}
          {field('Темп', 'pace', { placeholder: '5:15' })}
          {field('Пульс', 'hr', { placeholder: '150' })}
          {field('Набор, м', 'elevationGain', { placeholder: '40' })}
        </div>
        <button
          data-testid="save-activity-edit"
          className="btn"
          disabled={saving}
          onClick={() => void submit()}
          style={{ width: '100%', marginTop: 14, background: theme.accent, color: theme.onAccent, borderRadius: 14, padding: 14, opacity: saving ? 0.6 : 1 }}
        >
          Сохранить
        </button>
        <button className="btn" onClick={onClose} style={{ width: '100%', marginTop: 8, background: theme.card2, color: theme.dim, borderRadius: 14, padding: 12 }}>
          Отмена
        </button>
      </div>
    </div>,
    phoneEl,
  )
}
