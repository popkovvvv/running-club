import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ProfileAthlete } from './ProfileScreen'

const {
  joinClub,
  leaveClub,
  logout,
  stravaIntegration,
  connectStrava,
  disconnectStrava,
} = vi.hoisted(() => ({
  joinClub: vi.fn(),
  leaveClub: vi.fn(),
  logout: vi.fn(),
  stravaIntegration: vi.fn(),
  connectStrava: vi.fn(),
  disconnectStrava: vi.fn(),
}))

vi.mock('../lib/store', () => ({
  useApp: () => ({
    theme: {
      name: 'PULSE', display: 'Archivo', font: 'Manrope', radius: '20px',
      bg: '#0c0e10', card: '#16191d', card2: '#1e2228', line: '#2a2f36', text: '#f4f6f7',
      dim: '#9aa2ab', faint: '#5c636c', accent: '#ff5c22', accent2: '#ff7a45',
      accentSoft: 'rgba(255,92,34,.15)', good: '#39d98a', onAccent: '#ffffff',
    },
    user: { id: '1', name: 'Никита Попков', email: 'nikita@pulse.run', role: 'athlete', inClub: true },
    club: { id: 'c', name: 'PULSE', inviteCode: 'PULSE-7K42', accentHex: '#ff5c22', students: 1 },
    joinClub,
    leaveClub,
    logout,
  }),
}))

vi.mock('../lib/api', () => ({
  api: {
    stravaIntegration,
    connectStrava,
    disconnectStrava,
  },
}))

describe('ProfileAthlete Strava block', () => {
  beforeEach(() => {
    joinClub.mockReset()
    leaveClub.mockReset()
    logout.mockReset()
    stravaIntegration.mockReset()
    connectStrava.mockReset()
    disconnectStrava.mockReset()
  })

  it('renders disconnected Strava state', async () => {
    stravaIntegration.mockResolvedValue({
      provider: 'strava',
      status: 'disconnected',
      connected: false,
    })

    render(<ProfileAthlete />)

    await waitFor(() => expect(stravaIntegration).toHaveBeenCalled())
    expect(screen.getByTestId('strava-block')).toBeInTheDocument()
    expect(screen.getByText('Strava не подключена')).toBeInTheDocument()
    expect(screen.getByTestId('strava-connect')).toBeInTheDocument()
  })

  it('redirects to Strava connect url', async () => {
    const user = userEvent.setup()
    const assign = vi.fn()
    Object.defineProperty(window, 'location', {
      value: { assign },
      writable: true,
    })
    stravaIntegration.mockResolvedValue({
      provider: 'strava',
      status: 'disconnected',
      connected: false,
    })
    connectStrava.mockResolvedValue({ url: 'https://www.strava.com/oauth/authorize?x=1' })

    render(<ProfileAthlete />)

    await waitFor(() => expect(screen.getByTestId('strava-connect')).toBeInTheDocument())
    await user.click(screen.getByTestId('strava-connect'))

    expect(connectStrava).toHaveBeenCalled()
    expect(assign).toHaveBeenCalledWith('https://www.strava.com/oauth/authorize?x=1')
  })

  it('disconnects Strava and refreshes status', async () => {
    const user = userEvent.setup()
    stravaIntegration
      .mockResolvedValueOnce({
        provider: 'strava',
        status: 'active',
        connected: true,
        externalAthleteId: '424242',
      })
      .mockResolvedValueOnce({
        provider: 'strava',
        status: 'disconnected',
        connected: false,
      })
    disconnectStrava.mockResolvedValue({ status: 'ok' })

    render(<ProfileAthlete />)

    await waitFor(() => expect(screen.getByTestId('strava-disconnect')).toBeInTheDocument())
    await user.click(screen.getByTestId('strava-disconnect'))

    expect(disconnectStrava).toHaveBeenCalled()
    await waitFor(() => expect(screen.getByText('Strava не подключена')).toBeInTheDocument())
  })
})
