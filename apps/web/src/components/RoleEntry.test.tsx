import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { RoleEntry } from './RoleEntry'

const onPick = vi.fn()

vi.mock('../lib/store', async () => {
  const { pulseTheme } = await import('../lib/theme')
  return {
    useApp: () => ({
      theme: pulseTheme,
    }),
  }
})

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
