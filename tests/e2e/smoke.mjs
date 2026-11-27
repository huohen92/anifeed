import { chromium } from 'playwright';

const base = process.env.ANIFEED_URL || 'http://127.0.0.1:8115';
const errors = [];
const browser = await chromium.launch({ headless: true });
try {
  const page = await browser.newPage({ viewport: { width: 390, height: 844 } });
  page.on('pageerror', error => errors.push(error.message));
  await page.goto(base + '/', { waitUntil: 'domcontentloaded' });
  await page.waitForSelector('#view', { timeout: 15000 });
  const hasShell = await page.locator('.sidebar, .topbar').count();
  if (!hasShell) throw new Error('app shell did not render');
  const gate = await page.locator('#auth-gate').isVisible();
  if (gate) {
    if (!(await page.locator('#auth-user').count())) throw new Error('login form missing');
  } else {
    await page.waitForTimeout(1200);
    const view = (await page.locator('#view').innerText()).trim();
    if (!view) throw new Error('dashboard did not render');
  }
  if (!gate) {
    for (const route of ['feeds', 'tasks', 'history', 'logs', 'settings', 'ai']) {
      await page.evaluate(hash => { location.hash = hash; }, '#/' + route);
      await page.waitForTimeout(350);
      if (!(await page.locator('#view').innerText()).trim()) throw new Error(route + ' page is empty');
    }
  }
  if (errors.length) throw new Error('page errors: ' + errors.join('; '));
  console.log('e2e smoke passed, auth gate =', gate);
} finally {
  await browser.close();
}
