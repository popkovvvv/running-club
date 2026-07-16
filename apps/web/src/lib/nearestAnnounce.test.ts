import { describe, expect, it } from 'vitest'
import type { Announce } from './api'
import { nearestAnnounce } from './nearestAnnounce'

function ann(partial: Partial<Announce> & { id: string; place: string }): Announce {
  return {
    day: '',
    time: '19:50',
    group: 'Основная',
    note: '',
    going: 0,
    signedUp: false,
    scheduleCta: '',
    ...partial,
  }
}

function isoOffset(days: number): string {
  const d = new Date()
  d.setDate(d.getDate() + days)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

describe('nearestAnnounce', () => {
  it('picks closest upcoming by startsOn, not the last added', () => {
    const next = nearestAnnounce([
      ann({ id: 'far', place: 'Через неделю', startsOn: isoOffset(7) }),
      ann({ id: 'near', place: 'Завтра', startsOn: isoOffset(1) }),
      ann({ id: 'past', place: 'Вчера', startsOn: isoOffset(-1) }),
    ])
    expect(next?.id).toBe('near')
    expect(next?.place).toBe('Завтра')
  })

  it('falls back when all are in the past', () => {
    const next = nearestAnnounce([
      ann({ id: 'a', place: 'Старая', startsOn: isoOffset(-10) }),
      ann({ id: 'b', place: 'Ещё старее', startsOn: isoOffset(-20) }),
    ])
    expect(next?.id).toBe('a')
  })
})
