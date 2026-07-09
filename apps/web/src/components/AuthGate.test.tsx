import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AuthGate } from './AuthGate'
import { ProfileAthlete, ProfileCoach } from '../screens/ProfileScreen'

const register = vi.fn()
const login = vi.fn()
const onBack = vi.fn()

vi.mock('../lib/store', () => ({
  useApp: () => ({
    theme: {
      name: 'PULSE', display: 'Archivo', font: 'Manrope', radius: '20px',
      bg: '#0c0e10', card: '#16191d', card2: '#1e2228', line: '#2a2f36', text: '#f4f6f7',
      dim: '#9aa2ab', faint: '#5c636c', accent: '#ff5c22', accent2: '#ff7a45',
      accentSoft: 'rgba(255,92,34,.15)', good: '#39d98a', onAccent: '#ffffff',
    },
    login,
    register,
    user: { id: '1', name: 'Никита Попков', email: 'nikita@pulse.run', role: 'athlete', inClub: false },
    club: { id: 'c', name: 'PULSE', inviteCode: 'PULSE-7K42', accentHex: '#ff5c22', students: 1 },
    joinClub: vi.fn(),
    leaveClub: vi.fn(),
    setPalette: vi.fn(),
    logout: vi.fn(),
  }),
}))

describe('AuthGate', () => {
  beforeEach(() => {
    register.mockReset()
    login.mockReset()
    onBack.mockReset()
  })

  it('renders login form with empty credentials and role label', () => {
    render(<AuthGate role="athlete" onBack={onBack} />)
    expect(screen.getByTestId('auth-gate')).toBeInTheDocument()
    expect(screen.getByTestId('auth-submit')).toHaveTextContent('Войти')
    expect(screen.getByTestId('email-input')).toHaveValue('')
    expect(screen.getByTestId('password-input')).toHaveValue('')
    expect(screen.getByText(/Спортсмен/)).toBeInTheDocument()
    expect(screen.queryByText('ВЫ ВХОДИТЕ КАК')).toBeNull()
  })

  it('calls onBack when back is pressed', async () => {
    const user = userEvent.setup()
    render(<AuthGate role="coach" onBack={onBack} />)
    await user.click(screen.getByTestId('auth-back'))
    expect(onBack).toHaveBeenCalled()
  })

  it('registers with prop role', async () => {
    const user = userEvent.setup()
    render(<AuthGate role="coach" onBack={onBack} />)
    await user.click(screen.getByText('Регистрация'))
    await user.type(screen.getByTestId('name-input'), 'Иван Тренер')
    await user.type(screen.getByTestId('email-input'), 'coach@pulse.run')
    await user.type(screen.getByTestId('password-input'), 'secret')
    await user.click(screen.getByTestId('auth-submit'))
    expect(register).toHaveBeenCalledWith('Иван Тренер', 'coach@pulse.run', 'secret', 'coach')
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
