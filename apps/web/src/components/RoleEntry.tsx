import { useApp } from '../lib/store'

type Role = 'athlete' | 'coach'

export function RoleEntry({ onPick }: { onPick: (role: Role) => void }) {
  const { theme } = useApp()

  return (
    <div
      data-testid="role-entry"
      style={{
        position: 'absolute', inset: 0, zIndex: 50,
        background: `radial-gradient(120% 70% at 50% 0%,${theme.card} 0%,${theme.bg} 55%)`,
        display: 'flex', flexDirection: 'column', padding: '70px 26px 34px',
      }}
    >
      <div style={{ textAlign: 'center', marginBottom: 28 }}>
        <div style={{ fontFamily: theme.display, fontSize: 44, fontWeight: 900, color: theme.accent }}>{theme.name}</div>
        <div style={{ fontSize: 13, color: theme.dim, fontWeight: 600, marginTop: 4 }}>Клуб бегового прогресса</div>
      </div>
      <div className="card">
        <div style={{ fontFamily: theme.display, fontSize: 20, fontWeight: 800, marginBottom: 8 }}>Кто вы?</div>
        <div style={{ fontSize: 13, color: theme.dim, marginBottom: 18 }}>Выберите роль, чтобы войти или создать аккаунт</div>
        <button
          data-testid="pick-athlete"
          className="btn"
          onClick={() => onPick('athlete')}
          style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 13, padding: 15, fontSize: 15, marginBottom: 10 }}
        >
          Спортсмен
        </button>
        <button
          data-testid="pick-coach"
          className="btn"
          onClick={() => onPick('coach')}
          style={{ width: '100%', background: theme.card2, color: theme.text, borderRadius: 13, padding: 15, fontSize: 15, border: `1px solid ${theme.line}` }}
        >
          Тренер
        </button>
      </div>
    </div>
  )
}
