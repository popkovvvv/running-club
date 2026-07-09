import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { RoleEntry } from './RoleEntry'

const onPick = vi.fn()

vi.mock('../lib/store', () => ({
  useApp: () => ({
    theme: {
      name: 'PULSE', display: 'Archivo', font: 'Manrope', radius: '20px',
      bg: '#0c0e10', card: '#16191d', card2: '#1e2228', line: '#2a2f36', text: '#f4f6f7',
      dim: '#9aa2ab', faint: '#5c636c', accent: '#ff5c22', accent2: '#ff7a45',
      accentSoft: 'rgba(255,92,34,.15)', good: '#39d98a', onAccent: '#ffffff',
    },
  }),
}))

describe('RoleEntry', () => {
  beforeEach(() => {
    onPick.mockReset()
  })

  it('renders athlete and coach choices', () => {
    render(<RoleEntry onPick={onPick} />)
    expect(screen.getByTestId('role-entry')).toBeInTheDocument()
    expect(screen.getByTestId('pick-athlete')).toHaveTextContent('Спортсмен')
    expect(screen.getByTestId('pick-coach')).toHaveTextContent('Тренер')
  })

  it('calls onPick with athlete', async () => {
    const user = userEvent.setup()
    render(<RoleEntry onPick={onPick} />)
    await user.click(screen.getByTestId('pick-athlete'))
    expect(onPick).toHaveBeenCalledWith('athlete')
  })

  it('calls onPick with coach', async () => {
    const user = userEvent.setup()
    render(<RoleEntry onPick={onPick} />)
    await user.click(screen.getByTestId('pick-coach'))
    expect(onPick).toHaveBeenCalledWith('coach')
  })
})
