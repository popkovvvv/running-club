import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { CreateClubGate } from './CreateClubGate'

const createClub = vi.fn()

vi.mock('../lib/store', () => ({
  useApp: () => ({
    theme: {
      name: 'PULSE', display: 'Archivo', font: 'Manrope', radius: '20px',
      bg: '#0c0e10', card: '#16191d', card2: '#1e2228', line: '#2a2f36', text: '#f4f6f7',
      dim: '#9aa2ab', faint: '#5c636c', accent: '#ff5c22', accent2: '#ff7a45',
      accentSoft: 'rgba(255,92,34,.15)', good: '#39d98a', onAccent: '#ffffff',
    },
    createClub,
  }),
}))

describe('CreateClubGate', () => {
  beforeEach(() => {
    createClub.mockReset()
  })

  it('renders name and accent controls', () => {
    render(<CreateClubGate />)
    expect(screen.getByTestId('create-club-gate')).toBeInTheDocument()
    expect(screen.getByTestId('club-name-input')).toBeInTheDocument()
    expect(screen.getByTestId('create-club-submit')).toHaveTextContent('Создать клуб')
  })

  it('submits createClub with name and accent', async () => {
    const user = userEvent.setup()
    createClub.mockResolvedValue(undefined)
    render(<CreateClubGate />)
    await user.type(screen.getByTestId('club-name-input'), 'Night Run')
    await user.click(screen.getByTestId('accent-#4a9eff'))
    await user.click(screen.getByTestId('create-club-submit'))
    expect(createClub).toHaveBeenCalledWith({ name: 'Night Run', accentHex: '#4a9eff' })
  })

  it('shows error when name is empty', async () => {
    const user = userEvent.setup()
    render(<CreateClubGate />)
    await user.click(screen.getByTestId('create-club-submit'))
    expect(screen.getByTestId('create-club-error')).toHaveTextContent('Укажите название клуба')
    expect(createClub).not.toHaveBeenCalled()
  })
})
