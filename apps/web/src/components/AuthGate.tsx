import { useState } from 'react'
import { useApp } from '../lib/store'

export function AuthGate() {
  const { theme, login, register } = useApp()
  const [tab, setTab] = useState<'login' | 'register'>('login')
  const [role, setRole] = useState<'athlete' | 'coach'>('athlete')
  const [name, setName] = useState('')
  const [email, setEmail] = useState('nikita@pulse.run')
  const [password, setPassword] = useState('password')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  const submit = async () => {
    setBusy(true)
    setError('')
    try {
      if (tab === 'login') await login(email, password)
      else await register(name || 'Спортсмен', email, password, role)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div
      data-testid="auth-gate"
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
        <div style={{ display: 'flex', background: theme.card2, borderRadius: 12, padding: 4, gap: 4, marginBottom: 18 }}>
          <button className="btn" onClick={() => setTab('login')} style={{ flex: 1, padding: '10px 0', borderRadius: 9, background: tab === 'login' ? theme.accent : 'transparent', color: tab === 'login' ? theme.onAccent : theme.dim }}>Вход</button>
          <button className="btn" onClick={() => setTab('register')} style={{ flex: 1, padding: '10px 0', borderRadius: 9, background: tab === 'register' ? theme.accent : 'transparent', color: tab === 'register' ? theme.onAccent : theme.dim }}>Регистрация</button>
        </div>
        {tab === 'register' && (
          <input data-testid="name-input" placeholder="Имя и фамилия" value={name} onChange={(e) => setName(e.target.value)}
            style={{ width: '100%', background: theme.card2, border: `1px solid ${theme.line}`, borderRadius: 12, padding: 14, color: theme.text, marginBottom: 10, outline: 'none' }} />
        )}
        <input data-testid="email-input" placeholder="Логин или e-mail" value={email} onChange={(e) => setEmail(e.target.value)}
          style={{ width: '100%', background: theme.card2, border: `1px solid ${theme.line}`, borderRadius: 12, padding: 14, color: theme.text, marginBottom: 10, outline: 'none' }} />
        <input data-testid="password-input" type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)}
          style={{ width: '100%', background: theme.card2, border: `1px solid ${theme.line}`, borderRadius: 12, padding: 14, color: theme.text, marginBottom: 16, outline: 'none' }} />
        <div style={{ fontSize: 11, fontWeight: 800, color: theme.dim, letterSpacing: '.5px', marginBottom: 8 }}>ВЫ ВХОДИТЕ КАК</div>
        <div style={{ display: 'flex', background: theme.card2, borderRadius: 12, padding: 4, gap: 4, marginBottom: 20 }}>
          <button className="btn" onClick={() => setRole('athlete')} style={{ flex: 1, padding: '10px 0', borderRadius: 9, background: role === 'athlete' ? theme.accent : 'transparent', color: role === 'athlete' ? theme.onAccent : theme.dim }}>Спортсмен</button>
          <button className="btn" onClick={() => setRole('coach')} style={{ flex: 1, padding: '10px 0', borderRadius: 9, background: role === 'coach' ? theme.accent : 'transparent', color: role === 'coach' ? theme.onAccent : theme.dim }}>Тренер</button>
        </div>
        {error && <div data-testid="auth-error" style={{ color: theme.accent, fontSize: 13, marginBottom: 12 }}>{error}</div>}
        <button data-testid="auth-submit" className="btn" disabled={busy} onClick={() => void submit()}
          style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 13, padding: 15, fontSize: 15 }}>
          {tab === 'login' ? 'Войти' : 'Создать аккаунт'}
        </button>
      </div>
    </div>
  )
}
