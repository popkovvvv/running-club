import { useEffect, useState } from 'react'
import { api, type IntegrationStatus } from '../lib/api'
import { useApp } from '../lib/store'

const PALETTE = ['#ff5c22', '#c8ff34', '#4a9eff', '#ff3d81', '#22d3c5', '#ffd23f']

function initials(name?: string) {
  if (!name?.trim()) return '—'
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return parts[0].slice(0, 2).toUpperCase()
}

export function ProfileAthlete() {
  const { theme, user, club, joinClub, leaveClub, logout, tabEpoch } = useApp()
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [strava, setStrava] = useState<IntegrationStatus | null>(null)
  const [stravaLoading, setStravaLoading] = useState(true)
  const [stravaError, setStravaError] = useState('')

  useEffect(() => {
    let mounted = true
    setStravaLoading(true)
    void api.stravaIntegration()
      .then((res) => {
        if (mounted) {
          setStrava(res)
          setStravaError('')
        }
      })
      .catch((e) => {
        if (mounted) {
          setStravaError(e instanceof Error ? e.message : 'Ошибка')
        }
      })
      .finally(() => {
        if (mounted) {
          setStravaLoading(false)
        }
      })
    return () => {
      mounted = false
    }
  }, [tabEpoch])

  const connectStrava = async () => {
    const res = await api.connectStrava()
    window.location.assign(res.url)
  }

  const disconnectStrava = async () => {
    await api.disconnectStrava()
    const next = await api.stravaIntegration()
    setStrava(next)
    setStravaError('')
  }

  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }} data-testid="profile-athlete">
      <div className="card" style={{ textAlign: 'center' }}>
        <div style={{ width: 74, height: 74, borderRadius: '50%', margin: '0 auto', background: theme.card2, border: `2px solid ${theme.accent}`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontFamily: theme.display, fontWeight: 800, fontSize: 26, color: theme.accent }}>{initials(user?.name)}</div>
        <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800, marginTop: 12 }}>{user?.name}</div>
        <div style={{ fontSize: 12, color: theme.dim }}>{user?.email}</div>
      </div>
      <div className="card" data-testid="club-membership">
        <div style={{ fontFamily: theme.display, fontWeight: 800, marginBottom: 8 }}>Клуб</div>
        {user?.inClub ? (
          <>
            <div style={{ fontSize: 14 }}>Вы в клубе «{club?.name || '—'}»</div>
            <button data-testid="leave-club" className="btn" onClick={() => void leaveClub()} style={{ marginTop: 12, width: '100%', background: theme.card2, color: theme.accent, borderRadius: 12, padding: 12 }}>Покинуть клуб</button>
          </>
        ) : (
          <>
            <input data-testid="club-code" value={code} onChange={(e) => setCode(e.target.value)} placeholder="Код приглашения"
              style={{ width: '100%', background: theme.card2, border: `1px solid ${theme.line}`, borderRadius: 12, padding: 12, color: theme.text, outline: 'none' }} />
            {error && <div style={{ color: theme.accent, fontSize: 12, marginTop: 8 }}>{error}</div>}
            <button
              data-testid="join-club"
              className="btn"
              onClick={() => void joinClub(code).catch((e) => setError(e instanceof Error ? e.message : 'Ошибка'))}
              style={{ marginTop: 12, width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 12, padding: 12 }}
            >
              Вступить по коду
            </button>
          </>
        )}
      </div>
      <div className="card" data-testid="strava-block">
        <div style={{ fontFamily: theme.display, fontWeight: 800, marginBottom: 8 }}>Strava</div>
        {stravaLoading ? (
          <div style={{ fontSize: 13, color: theme.dim }}>Проверяем подключение…</div>
        ) : strava?.connected ? (
          <>
            <div style={{ fontSize: 14 }}>Подключена</div>
            <div style={{ fontSize: 12, color: theme.dim, marginTop: 6 }}>Athlete ID: {strava.externalAthleteId || '—'}</div>
            <div style={{ fontSize: 12, color: theme.dim, marginTop: 4 }}>Последняя синхронизация: {strava.lastSyncedAt || 'ещё не было'}</div>
            <button
              data-testid="strava-disconnect"
              className="btn"
              onClick={() => void disconnectStrava()}
              style={{ marginTop: 12, width: '100%', background: theme.card2, color: theme.accent, borderRadius: 12, padding: 12 }}
            >
              Отключить Strava
            </button>
          </>
        ) : (
          <>
            <div style={{ fontSize: 14 }}>Strava не подключена</div>
            {stravaError && <div style={{ color: theme.accent, fontSize: 12, marginTop: 8 }}>{stravaError}</div>}
            <button
              data-testid="strava-connect"
              className="btn"
              onClick={() => void connectStrava()}
              style={{ marginTop: 12, width: '100%', background: theme.accent, color: theme.onAccent, borderRadius: 12, padding: 12 }}
            >
              Подключить Strava
            </button>
          </>
        )}
      </div>
      <button className="btn" onClick={logout} style={{ background: theme.card2, color: theme.dim, borderRadius: 14, padding: 14 }}>Выйти</button>
    </div>
  )
}

export function ProfileCoach() {
  const { theme, club, setPalette, logout } = useApp()
  return (
    <div className="fade" style={{ display: 'flex', flexDirection: 'column', gap: 14 }} data-testid="profile-coach">
      <div className="card" style={{ textAlign: 'center' }}>
        <div style={{ fontFamily: theme.display, fontSize: 22, fontWeight: 800 }}>Главный тренер</div>
        <div style={{ fontSize: 12, color: theme.dim }}>Клуб «{club?.name || '—'}» · {club?.students ?? 0} учеников</div>
      </div>
      <div className="card" data-testid="invite-code-block">
        <div style={{ fontFamily: theme.display, fontWeight: 800 }}>Код приглашения</div>
        <div data-testid="invite-code" style={{ fontFamily: theme.display, fontSize: 28, fontWeight: 800, color: theme.accent, marginTop: 8 }}>{club?.inviteCode || '—'}</div>
      </div>
      <div className="card" data-testid="palette-block">
        <div style={{ fontFamily: theme.display, fontWeight: 800 }}>Палитра клуба</div>
        <div style={{ fontSize: 12, color: theme.dim, margin: '6px 0 16px' }}>Акцент применяется ко всем участникам.</div>
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
          {PALETTE.map((hex) => (
            <button
              key={hex}
              data-testid={`palette-${hex}`}
              className="btn"
              onClick={() => void setPalette(hex)}
              style={{ width: 46, height: 46, borderRadius: 14, background: hex, border: `3px solid ${club?.accentHex === hex ? theme.text : 'transparent'}` }}
            />
          ))}
        </div>
      </div>
      <button className="btn" onClick={logout} style={{ background: theme.card2, color: theme.dim, borderRadius: 14, padding: 14 }}>Выйти</button>
    </div>
  )
}
