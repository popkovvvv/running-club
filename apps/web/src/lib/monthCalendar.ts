import type { ScheduleCalCell } from '../components/schedule/ScheduleCalendar'

export const MONTH_LABELS = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь',
]

export function buildMonthCells(
  year: number,
  month: number,
  theme: { text: string; accent: string; accentSoft: string; good: string },
  marked: Set<string>,
  completed: Set<string>,
): ScheduleCalCell[] {
  const first = new Date(year, month - 1, 1)
  const daysInMonth = new Date(year, month, 0).getDate()
  const startDow = (first.getDay() + 6) % 7
  const now = new Date()
  const isCurrentMonth = year === now.getFullYear() && month === now.getMonth() + 1
  const today = now.getDate()
  const cells: ScheduleCalCell[] = []

  for (let i = 0; i < startDow; i++) {
    cells.push({ key: `b${i}`, n: 0, blank: true, bg: 'transparent', fg: theme.text, dot: 'transparent' })
  }
  for (let d = 1; d <= daysInMonth; d++) {
    const iso = `${year}-${String(month).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const has = marked.has(iso)
    const done = completed.has(iso)
    const isToday = isCurrentMonth && d === today
    let bg = 'transparent'
    const fg = theme.text
    let dot = 'transparent'
    if (isToday) bg = theme.accentSoft
    if (has) {
      dot = done ? theme.good : theme.accent
      if (!isToday) bg = theme.accentSoft
    }
    cells.push({ key: `d${d}`, n: d, blank: false, has, isToday, iso, bg, fg, dot })
  }
  return cells
}
