import { test, expect } from '@playwright/test'

async function pickRoleAndLogin(page: import('@playwright/test').Page, role: 'athlete' | 'coach', email: string) {
  await expect(page.getByTestId('role-entry')).toBeVisible()
  await page.getByTestId(role === 'coach' ? 'pick-coach' : 'pick-athlete').click()
  await expect(page.getByTestId('auth-gate')).toBeVisible()
  await page.getByTestId('email-input').fill(email)
  await page.getByTestId('password-input').fill('password')
  await page.getByTestId('auth-submit').click()
}

test.beforeEach(async ({ page }) => {
  await page.goto('/')
  await page.evaluate(() => localStorage.clear())
  await page.reload()
  await expect(page.getByTestId('role-entry')).toBeVisible()
})

test.describe('running club', () => {
  test('athlete login, schedule cta, no payments', async ({ page }) => {
    await pickRoleAndLogin(page, 'athlete', 'nikita@pulse.run')
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await expect(page.getByText('Личный кабинет')).toBeVisible()
    await expect(page.getByTestId('nav-schedule')).toHaveCount(0)
    await expect(page.getByTestId('nav-home')).toBeVisible()
    await expect(page.getByTestId('nav-plan')).toBeVisible()
    await expect(page.getByTestId('nav-prog')).toBeVisible()
    await expect(page.getByTestId('nav-races')).toBeVisible()
    await expect(page.getByTestId('nav-profile')).toBeVisible()

    await page.getByTestId('group-preview').click()
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
    await pickRoleAndLogin(page, 'coach', 'coach@pulse.run')
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await page.getByTestId('nav-profile').click()
    await expect(page.getByTestId('invite-code')).toContainText('PULSE')
    await page.getByTestId('palette-#c8ff34').click()
    await expect(page.getByText('Оплаты учеников')).toHaveCount(0)
  })

  test('calendar blank cells', async ({ page }) => {
    await pickRoleAndLogin(page, 'athlete', 'nikita@pulse.run')
    await expect(page.getByTestId('bottom-nav')).toBeVisible({ timeout: 20000 })
    await page.getByTestId('group-preview').click()
    await expect(page.getByTestId('cal-blank').first()).toBeVisible()
  })
})
