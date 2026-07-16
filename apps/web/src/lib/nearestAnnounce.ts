import type { Announce } from '../lib/api'

function todayIso(): string {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

export function nearestAnnounce(announces: Announce[]): Announce | undefined {
  if (announces.length === 0) return undefined
  const today = todayIso()
  const upcoming = announces
    .filter((a) => !a.startsOn || a.startsOn >= today)
    .sort((a, b) => {
      if (a.startsOn && b.startsOn) return a.startsOn.localeCompare(b.startsOn)
      if (a.startsOn) return -1
      if (b.startsOn) return 1
      return 0
    })
  return upcoming[0] ?? announces[0]
}
