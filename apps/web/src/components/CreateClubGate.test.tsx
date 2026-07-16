import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { CreateClubGate } from './CreateClubGate'

const createClub = vi.hoisted(() => vi.fn())

vi.mock('../lib/store', async () => {
  const { pulseTheme } = await import('../lib/theme')
  return {
    useApp: () => ({
      theme: pulseTheme,
      createClub,
    }),
  }
})

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
