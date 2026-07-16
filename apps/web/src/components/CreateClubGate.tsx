import { useState } from 'react'
import { useApp } from '../lib/store'

const PALETTE = ['#ff5c22', '#c8ff34', '#4a9eff', '#ff3d81', '#22d3c5', '#ffd23f']

export function CreateClubGate() {
  const { theme, createClub } = useApp()
  const [name, setName] = useState('')
  const [accentHex, setAccentHex] = useState(PALETTE[0])
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  const submit = async () => {
    if (!name.trim()) {
      setError('Укажите название клуба')
      return
    }
    setBusy(true)
    setError('')
    try {
      await createClub({ name: name.trim(), accentHex })
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div
      data-testid="create-club-gate"
      className="gate-atmos"
      style={{
        position: 'absolute', inset: 0, zIndex: 50,
        display: 'flex', flexDirection: 'column', padding: '70px 26px 34px',
      }}
    >
      <div style={{ textAlign: 'center', marginBottom: 28 }}>
        <div className="display-title" style={{ fontSize: 40, color: theme.accent }}>{name.trim() || theme.name}</div>
        <div style={{ fontSize: 13, color: theme.dim, fontWeight: 600, marginTop: 8 }}>Создайте свой клуб</div>
      </div>
      <div className="card">
        <div style={{ fontFamily: theme.display, fontSize: 20, fontWeight: 800, marginBottom: 16 }}>Новый клуб</div>
        <input
          data-testid="club-name-input"
          placeholder="Название клуба"
          value={name}
          onChange={(e) => setName(e.target.value)}
          style={{ width: '100%', background: theme.card2, border: `1px solid ${theme.line}`, borderRadius: 12, padding: 14, color: theme.text, marginBottom: 16, outline: 'none' }}
        />
        <div style={{ fontSize: 11, fontWeight: 800, color: theme.dim, letterSpacing: '.5px', marginBottom: 10 }}>АКЦЕНТ</div>
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', marginBottom: 20 }}>
          {PALETTE.map((hex) => (
            <button
              key={hex}
              data-testid={`accent-${hex}`}
              className="btn"
              onClick={() => setAccentHex(hex)}
              style={{ width: 46, height: 46, borderRadius: 14, background: hex, border: `3px solid ${accentHex === hex ? theme.text : 'transparent'}` }}
            />
          ))}
        </div>
        {error && <div data-testid="create-club-error" style={{ color: theme.accent, fontSize: 13, marginBottom: 12 }}>{error}</div>}
        <button
          data-testid="create-club-submit"
          className="btn"
          disabled={busy}
          onClick={() => void submit()}
          style={{ width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 13, padding: 15, fontSize: 15 }}
        >
          Создать клуб
        </button>
      </div>
    </div>
  )
}
