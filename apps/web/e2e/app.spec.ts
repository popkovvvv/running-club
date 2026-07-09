import { test, expect } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.goto('/')
  await page.evaluate(() => localStorage.clear())
  await page.reload()
  await expect(page.getByTestId('auth-gate')).toBeVisible()
})

test.describe('running club', () => {
  test('athlete login, schedule cta, no payments', async ({ page }) => {
    await page.getByTestId('email-input').fill('nikita@pulse.run')
    await page.getByTestId('password-input').fill('password')
    await page.getByTestId('auth-submit').click()
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await expect(page.getByText('Личный кабинет')).toBeVisible()

    await page.getByTestId('nav-schedule').click()
    const cta = page.getByTestId('schedule-cta').first()
    await expect(cta).toHaveText(/Записаться|Вы записаны/)
    if ((await cta.textContent()) === 'Записаться') {
      await cta.click()
      await expect(cta).toHaveText('Вы записаны')
    }

    await page.getByTestId('nav-profile').click()
    await expect(page.getByTestId('profile-athlete')).toBeVisible()
    await expect(page.getByText('Абонемент')).toHaveCount(0)
    await expect(page.getByText('Оплаты')).toHaveCount(0)
  })

  test('coach palette and invite code', async ({ page }) => {
    await page.getByTestId('email-input').fill('coach@pulse.run')
    await page.getByTestId('password-input').fill('password')
    await page.getByTestId('auth-submit').click()
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await page.getByTestId('nav-profile').click()
    await expect(page.getByTestId('invite-code')).toContainText('PULSE')
    await page.getByTestId('palette-#c8ff34').click()
    await expect(page.getByText('Оплаты учеников')).toHaveCount(0)
  })

  test('calendar blank cells', async ({ page }) => {
    await page.getByTestId('email-input').fill('nikita@pulse.run')
    await page.getByTestId('password-input').fill('password')
    await page.getByTestId('auth-submit').click()
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await page.getByTestId('nav-schedule').click()
    await expect(page.getByTestId('cal-blank').first()).toBeVisible()
  })
})
