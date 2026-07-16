import { useApp } from '../lib/store'

type Role = 'athlete' | 'coach'

export function RoleEntry({ onPick }: { onPick: (role: Role) => void }) {
  const { theme } = useApp()

  return (
    <div
      data-testid="role-entry"
      className="gate-atmos gate-shell"
    >
      <div style={{ textAlign: 'center', marginBottom: 28 }}>
        <div className="display-title" style={{ fontSize: 'clamp(32px, 8vw, 44px)', color: theme.accent }}>{theme.name}</div>
        <div style={{ fontSize: 13, color: theme.dim, fontWeight: 600, marginTop: 8 }}>Клуб бегового прогресса</div>
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
