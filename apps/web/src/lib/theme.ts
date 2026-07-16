export type Theme = {
  name: string
  font: string
  display: string
  radius: string
  bg: string
  card: string
  card2: string
  line: string
  text: string
  dim: string
  faint: string
  accent: string
  accent2: string
  accentSoft: string
  good: string
  onAccent: string
  heroGlow: string
  halftoneOpacity: string
  heroImage: string
}

export const pulseTheme: Theme = {
  name: 'PULSE',
  font: 'Manrope',
  display: 'Archivo',
  radius: '20px',
  bg: '#070b12',
  card: '#101820',
  card2: '#16202c',
  line: '#243040',
  text: '#f4f6f7',
  dim: '#8b96a5',
  faint: '#5a6573',
  accent: '#ff5c22',
  accent2: '#ff7a45',
  accentSoft: 'rgba(255,92,34,.15)',
  good: '#39d98a',
  onAccent: '#ffffff',
  heroGlow: 'rgba(0, 120, 220, 0.35)',
  halftoneOpacity: '0.22',
  heroImage: '',
}

function hexToRgb(hex: string): [number, number, number] {
  const h = hex.replace('#', '')
  return [parseInt(h.slice(0, 2), 16), parseInt(h.slice(2, 4), 16), parseInt(h.slice(4, 6), 16)]
}

export function withAccent(base: Theme, accentHex: string): Theme {
  const [ar, ag, ab] = hexToRgb(accentHex)
  const lum = (0.299 * ar + 0.587 * ag + 0.114 * ab) / 255
  const li = (c: number) => Math.round(c + (255 - c) * 0.25)
  return {
    ...base,
    accent: accentHex,
    accent2: `rgb(${li(ar)},${li(ag)},${li(ab)})`,
    accentSoft: `rgba(${ar},${ag},${ab},.15)`,
    onAccent: lum > 0.6 ? '#0a0b0c' : '#ffffff',
  }
}

export function themeToCssVars(t: Theme): Record<string, string> {
  return {
    '--bg': t.bg,
    '--card': t.card,
    '--card2': t.card2,
    '--line': t.line,
    '--text': t.text,
    '--dim': t.dim,
    '--faint': t.faint,
    '--accent': t.accent,
    '--accent2': t.accent2,
    '--accent-soft': t.accentSoft,
    '--good': t.good,
    '--on-accent': t.onAccent,
    '--radius': t.radius,
    '--font': t.font,
    '--display': t.display,
    '--hero-glow': t.heroGlow,
    '--halftone-opacity': t.halftoneOpacity,
    '--hero-image': t.heroImage ? `url(${t.heroImage})` : 'none',
  }
}

export function scheduleCta(signedUp: boolean): string {
  return signedUp ? 'Вы записаны' : 'Записаться'
}

export type CalendarCell = {
  key: string
  n: number
  blank: boolean
  has: boolean
  isToday: boolean
  bg: string
  fg: string
  dot: string
}

export function mapCalendarCells(
  cells: Array<{ key: string; n: number; blank: boolean; has: boolean; isToday: boolean }>,
  t: Theme,
): CalendarCell[] {
  return cells.map((c) => {
    if (c.blank) {
      return { ...c, bg: 'transparent', fg: t.dim, dot: 'transparent' }
    }
    return {
      ...c,
      bg: c.isToday ? t.accent : c.has ? t.accentSoft : 'transparent',
      fg: c.isToday ? t.onAccent : t.text,
      dot: c.has && !c.isToday ? t.accent : 'transparent',
    }
  })
}
