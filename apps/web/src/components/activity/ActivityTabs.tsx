import type { Theme } from '../../lib/theme'

export function ActivityTabs({ theme, tab, tabs, onChange }: {
  theme: Theme
  tab: string
  tabs: { id: string; label: string }[]
  onChange: (id: string) => void
}) {
  return (
    <div style={{ display: 'flex', gap: 6, overflowX: 'auto', paddingBottom: 4 }}>
      {tabs.map((t) => (
        <button
          key={t.id}
          type="button"
          className="btn"
          onClick={() => onChange(t.id)}
          style={{
            padding: '8px 14px',
            borderRadius: 20,
            fontSize: 12,
            fontWeight: 700,
            whiteSpace: 'nowrap',
            background: tab === t.id ? theme.accent : theme.card2,
            color: tab === t.id ? theme.onAccent : theme.dim,
          }}
        >
          {t.label}
        </button>
      ))}
    </div>
  )
}
