import type { Theme } from '../../lib/theme'

export function MetricGrid({ theme, items }: { theme: Theme; items: { label: string; value: string }[] }) {
  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 10 }}>
      {items.map((item) => (
        <div key={item.label} className="card" style={{ textAlign: 'center', padding: 12 }}>
          <div style={{ fontFamily: theme.display, fontWeight: 800, fontSize: 18 }}>{item.value}</div>
          <div style={{ fontSize: 10, color: theme.dim }}>{item.label}</div>
        </div>
      ))}
    </div>
  )
}
