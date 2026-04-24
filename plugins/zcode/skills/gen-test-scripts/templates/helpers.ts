import { execSync } from 'node:child_process';
import { writeFileSync, readFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parse as parseYaml } from 'yaml';
import { chromium, type Browser, type Page } from 'playwright';

const __dirname = dirname(fileURLToPath(import.meta.url));
const SCREENSHOTS_DIR = join(__dirname, '..', 'results', 'screenshots');

// ── Config ─────────────────────────────────────────────────────────
function readConfig(): Record<string, any> {
  const configPath = resolve('tests/e2e/config.yaml');
  if (!existsSync(configPath)) {
    throw new Error(`Config not found: ${configPath}. Create it with required fields.`);
  }
  return parseYaml(readFileSync(configPath, 'utf-8'));
}

const _config = readConfig();

export const baseUrl = _config.baseUrl ?? 'http://localhost:3456';
export const apiUrl = _config.apiUrl ?? 'http://localhost:8080';
const DEFAULT_TIMEOUT = parseInt(_config.timeout ?? '30000');

// ── Browser lifecycle ──────────────────────────────────────────────
let _browser: Browser | null = null;
let _page: Page | null = null;

export async function setupBrowser(): Promise<Page> {
  _browser = await chromium.launch();
  _page = await _browser.newPage();
  _page.setDefaultTimeout(DEFAULT_TIMEOUT);
  return _page;
}

export async function teardownBrowser(): Promise<void> {
  await _browser?.close();
  _browser = null;
  _page = null;
}

export function getPage(): Page {
  if (!_page) throw new Error('Browser not initialized. Call setupBrowser() first.');
  return _page;
}

// ── Evidence ───────────────────────────────────────────────────────
export async function screenshot(page: Page, tcId: string): Promise<string> {
  if (!existsSync(SCREENSHOTS_DIR)) mkdirSync(SCREENSHOTS_DIR, { recursive: true });
  const path = join(SCREENSHOTS_DIR, `${tcId}.png`);
  await page.screenshot({ path, fullPage: true });
  return path;
}

// ── HTTP ───────────────────────────────────────────────────────────
export interface CurlResponse {
  status: number;
  headers: Record<string, string>;
  body: string;
}

export async function curl(
  method: string,
  url: string,
  opts?: {
    body?: string;
    headers?: Record<string, string>;
    timeout?: number;
  },
): Promise<CurlResponse> {
  const controller = new AbortController();
  const timeout = setTimeout(
    () => controller.abort(),
    opts?.timeout ?? 10000,
  );

  try {
    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...opts?.headers,
      },
      body: opts?.body,
      signal: controller.signal,
    });

    const headers: Record<string, string> = {};
    res.headers.forEach((v, k) => { headers[k] = v; });

    return {
      status: res.status,
      headers,
      body: await res.text(),
    };
  } finally {
    clearTimeout(timeout);
  }
}

// ── Auth ────────────────────────────────────────────────────────────
export interface UICredentials {
  username: string;
  password: string;
}

export const defaultCreds: UICredentials = {
  username: _config.username ?? 'admin',
  password: _config.password ?? 'password',
};

export async function loginViaUI(page: Page, creds: UICredentials = defaultCreds): Promise<void> {
  await page.goto(`${baseUrl}/login`);
  await page.waitForLoadState('networkidle');
  await page.getByRole('textbox', { name: /username|email/i }).fill(creds.username);
  await page.getByRole('textbox', { name: /password/i }).fill(creds.password);
  await page.getByRole('button', { name: /login|sign in|submit/i }).click();
  await page.waitForURL((url) => !url.pathname.includes('login'), { timeout: DEFAULT_TIMEOUT });
}

export async function getApiToken(apiUrl: string, creds: UICredentials = defaultCreds): Promise<string> {
  const res = await curl('POST', `${apiUrl}/api/auth/login`, {
    body: JSON.stringify({ username: creds.username, password: creds.password }),
  });
  if (res.status !== 200) throw new Error(`Auth failed: ${res.status} ${res.body}`);
  const data = JSON.parse(res.body);
  return data.token ?? data.access_token ?? data.data?.token ?? '';
}

export function createAuthCurl(
  apiUrl: string,
  token: string,
): (method: string, path: string, opts?: { body?: string; headers?: Record<string, string>; timeout?: number }) => Promise<CurlResponse> {
  return (method, path, opts) =>
    curl(method, `${apiUrl}${path}`, {
      ...opts,
      headers: { Authorization: `Bearer ${token}`, ...opts?.headers },
    });
}

// ── CLI ────────────────────────────────────────────────────────────
export interface CliResult {
  stdout: string;
  stderr: string;
  exitCode: number;
}

export function runCli(cmd: string, cwd?: string): CliResult {
  try {
    const stdout = execSync(cmd, {
      encoding: 'utf-8',
      timeout: DEFAULT_TIMEOUT,
      cwd: cwd ?? process.cwd(),
    });
    return { stdout, stderr: '', exitCode: 0 };
  } catch (e: any) {
    return {
      stdout: e.stdout ?? '',
      stderr: e.stderr ?? '',
      exitCode: e.status ?? 1,
    };
  }
}
