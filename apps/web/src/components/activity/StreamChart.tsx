function parseStreamData(data: string): number[] {
  try {
    const parsed = JSON.parse(data)
    return Array.isArray(parsed) ? parsed.map(Number) : []
  } catch {
    return []
  }
}

export function StreamChart({ values, color, height = 100, label }: { values: number[]; color: string; height?: number; label?: string }) {
  if (values.length < 2) {
    return <div style={{ height, fontSize: 12, color: '#9aa2ab', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>Нет данных</div>
  }
  const min = Math.min(...values)
  const max = Math.max(...values)
  const span = max - min || 1
  const w = 300
  const points = values.map((v, i) => {
    const x = (i / (values.length - 1)) * w
    const y = height - ((v - min) / span) * (height - 8) - 4
    return `${x},${y}`
  }).join(' ')

  return (
    <div>
      {label && <div style={{ fontSize: 11, color: '#9aa2ab', marginBottom: 6 }}>{label}</div>}
      <svg viewBox={`0 0 ${w} ${height}`} style={{ width: '100%', height }}>
        <polyline points={points} fill="none" stroke={color} strokeWidth="2.5" strokeLinejoin="round" />
      </svg>
    </div>
  )
}

export function parseStreamValues(data: string): number[] {
  return parseStreamData(data)
}

export function computeSplits(latlngData: string, timeData: string): { km: number; pace: string; elev: string }[] {
  let coords: number[][] = []
  let times: number[] = []
  try {
    coords = JSON.parse(latlngData)
    times = JSON.parse(timeData)
  } catch {
    return []
  }
  if (coords.length < 2 || times.length < 2) return []

  const splits: { km: number; pace: string; elev: string }[] = []
  let kmAcc = 0
  let segStart = 0

  const dist = (a: number[], b: number[]) => {
    const R = 6371
    const dLat = ((b[0] - a[0]) * Math.PI) / 180
    const dLng = ((b[1] - a[1]) * Math.PI) / 180
    const x = Math.sin(dLat / 2) ** 2 + Math.cos((a[0] * Math.PI) / 180) * Math.cos((b[0] * Math.PI) / 180) * Math.sin(dLng / 2) ** 2
    return R * 2 * Math.atan2(Math.sqrt(x), Math.sqrt(1 - x))
  }

  for (let i = 1; i < coords.length; i++) {
    kmAcc += dist(coords[i - 1], coords[i])
    if (kmAcc >= splits.length + 1 || i === coords.length - 1) {
      const dt = times[i] - times[segStart]
      const paceSec = dt > 0 ? dt / Math.max(kmAcc - splits.length, 0.01) : 0
      const m = Math.floor(paceSec / 60)
      const s = Math.floor(paceSec % 60)
      splits.push({ km: splits.length + 1, pace: `${m}:${String(s).padStart(2, '0')}`, elev: '—' })
      segStart = i
    }
  }
  return splits.slice(0, 20)
}

export function SplitTable({ splits }: { splits: { km: number; pace: string; elev: string }[] }) {
  if (splits.length === 0) {
    return <div style={{ fontSize: 13, color: '#9aa2ab' }}>Недостаточно GPS-данных для сплитов</div>
  }
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
      {splits.map((s) => (
        <div key={s.km} style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13, padding: '8px 0', borderBottom: '1px solid rgba(255,255,255,.06)' }}>
          <span>Км {s.km}</span>
          <span style={{ fontWeight: 700 }}>{s.pace} /км</span>
        </div>
      ))}
    </div>
  )
}
