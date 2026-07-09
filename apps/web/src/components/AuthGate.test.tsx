import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { AuthGate } from '../components/AuthGate'
import { ProfileAthlete, ProfileCoach } from '../screens/ProfileScreen'

vi.mock('../lib/store', () => ({
  useApp: () => ({
    theme: {
      name: 'PULSE', display: 'Archivo', font: 'Manrope', radius: '20px',
      bg: '#0c0e10', card: '#16191d', card2: '#1e2228', line: '#2a2f36', text: '#f4f6f7',
      dim: '#9aa2ab', faint: '#5c636c', accent: '#ff5c22', accent2: '#ff7a45',
      accentSoft: 'rgba(255,92,34,.15)', good: '#39d98a', onAccent: '#ffffff',
    },
    login: vi.fn(),
    register: vi.fn(),
    user: { id: '1', name: 'Никита Попков', email: 'nikita@pulse.run', role: 'athlete', inClub: false },
    club: { id: 'c', name: 'PULSE', inviteCode: 'PULSE-7K42', accentHex: '#ff5c22', students: 1 },
    joinClub: vi.fn(),
    leaveClub: vi.fn(),
    setPalette: vi.fn(),
    logout: vi.fn(),
  }),
}))

describe('AuthGate', () => {
  it('renders login form', () => {
    render(<AuthGate />)
    expect(screen.getByTestId('auth-gate')).toBeInTheDocument()
    expect(screen.getByTestId('auth-submit')).toHaveTextContent('Войти')
  })
})

describe('profiles without payment', () => {
  it('athlete profile has no payment rows', () => {
    render(<ProfileAthlete />)
    expect(screen.queryByText(/Абонемент/i)).toBeNull()
    expect(screen.queryByText(/Оплат/i)).toBeNull()
    expect(screen.getByTestId('club-membership')).toBeInTheDocument()
  })

  it('coach profile has invite and palette, no payments', () => {
    render(<ProfileCoach />)
    expect(screen.getByTestId('invite-code')).toHaveTextContent('PULSE-7K42')
    expect(screen.getByTestId('palette-block')).toBeInTheDocument()
    expect(screen.queryByText(/Оплаты учеников/i)).toBeNull()
  })
})
