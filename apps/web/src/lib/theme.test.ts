import { describe, expect, it } from 'vitest'
import { mapCalendarCells, pulseTheme, scheduleCta, themeToCssVars, withAccent } from './theme'

describe('scheduleCta', () => {
  it.each([
    [false, 'Записаться'],
    [true, 'Вы записаны'],
  ] as const)('signedUp=%s → %s', (signed, label) => {
    expect(scheduleCta(signed)).toBe(label)
  })
})

describe('withAccent', () => {
  it.each([
    ['#ff5c22', '#ffffff'],
    ['#c8ff34', '#0a0b0c'],
  ] as const)('accent %s sets onAccent %s', (hex, onAccent) => {
    expect(withAccent(pulseTheme, hex).onAccent).toBe(onAccent)
    expect(withAccent(pulseTheme, hex).accent).toBe(hex)
  })
})

describe('themeToCssVars', () => {
  it('exposes hero tokens', () => {
    const vars = themeToCssVars(pulseTheme)
    expect(vars['--bg']).toBe('#070b12')
    expect(vars['--hero-glow']).toBe(pulseTheme.heroGlow)
    expect(vars['--halftone-opacity']).toBe(pulseTheme.halftoneOpacity)
    expect(vars['--hero-image']).toBe('none')
  })
})

describe('mapCalendarCells', () => {
  it('blank cells have transparent bg and dot', () => {
    const cells = mapCalendarCells(
      [
        { key: 'b0', n: 0, blank: true, has: false, isToday: false },
        { key: 'd9', n: 9, blank: false, has: true, isToday: true },
        { key: 'd21', n: 21, blank: false, has: true, isToday: false },
      ],
      pulseTheme,
    )
    expect(cells[0].blank).toBe(true)
    expect(cells[0].bg).toBe('transparent')
    expect(cells[0].dot).toBe('transparent')
    expect(cells[1].bg).toBe(pulseTheme.accent)
    expect(cells[2].dot).toBe(pulseTheme.accent)
  })
})
