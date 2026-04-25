import { execSync } from 'node:child_process';
import { readFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parse as parseYaml } from 'yaml';
import { chromium, type Browser, type Page } from 'playwright';

const __dirname = dirname(fileURLToPath(import.meta.url));

// ── Config ─────────────────────────────────────────────────────────
const _configPath = findConfigPath();

function findConfigPath(): string {
  // Allow explicit override via environment variable
  const envPath = process.env.E2E_CONFIG_PATH;
  if (envPath && existsSync(envPath)) return resolve(envPath);

  let dir = __dirname;
  for (let i = 0; i < 10; i++) {
    const candidate = resolve(dir, 'tests', 'e2e', 'config.yaml');
    if (existsSync(candidate)) return candidate;
    const parent = resolve(dir, '..');
    if (parent === dir) break;
    dir = parent;
  }
  throw new Error(`tests/e2e/config.yaml not found. Searched upward from ${__dirname}. Set E2E_CONFIG_PATH or run /gen-sitemap first.`);
}

// Screenshots go to <helpers-dir>/../results/screenshots
// During development: testing/scripts/ → testing/results/screenshots/
// After graduation:   tests/e2e/       → tests/e2e/results/screenshots/
const SCREENSHOTS_DIR = join(__dirname, '..', 'results', 'screenshots');

interface E2EConfig {
  baseUrl?: string;
  apiUrl?: string;
  timeout?: number | string;
  username?: string;
  password?: string;
  loginLocators?: { usernameField?: string; passwordField?: string; submitButton?: string };
}

function readConfig(): E2EConfig {
  return parseYaml(readFileSync(findConfigPath(), 'utf-8'));
}

const _config = readConfig();

function toNumber(val: unknown, fallback: number): number {
  if (typeof val === 'number' && Number.isFinite(val)) return val;
  if (typeof val === 'string') {
    const n = parseInt(val, 10);
    return Number.isFinite(n) ? n : fallback;
  }
  return fallback;
}

export const baseUrl = _config.baseUrl ?? 'http://localhost:3456';
export const apiUrl = _config.apiUrl ?? 'http://localhost:8080';
const DEFAULT_TIMEOUT = toNumber(_config.timeout, 30000);

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

const _loginLocators = _config.loginLocators;

// TEMPLATE: Replace regex with actual locator from sitemap when generating
export async function loginViaUI(page: Page, creds: UICredentials = defaultCreds): Promise<void> {
  await page.goto(`${baseUrl}/login`);
  await page.waitForLoadState('networkidle');
  const uPat = new RegExp(_loginLocators?.usernameField ?? 'username|email', 'i');
  const pPat = new RegExp(_loginLocators?.passwordField ?? 'password', 'i');
  const bPat = new RegExp(_loginLocators?.submitButton ?? 'login|sign in|submit', 'i');
  await page.getByRole('textbox', { name: uPat }).fill(creds.username);
  await page.getByRole('textbox', { name: pPat }).fill(creds.password);
  await page.getByRole('button', { name: bPat }).click();
  await page.waitForURL((url) => !url.pathname.includes('login'), { timeout: DEFAULT_TIMEOUT });
}

export async function getApiToken(apiUrl: string, creds: UICredentials = defaultCreds): Promise<string> {
  const res = await curl('POST', `${apiUrl}/api/auth/login`, {
    body: JSON.stringify({ username: creds.username, password: creds.password }),
  });
  if (res.status !== 200) throw new Error(`Auth failed: ${res.status} ${res.body}`);
  const data = JSON.parse(res.body);
  const token = data.token ?? data.access_token ?? data.data?.token;
  if (!token) throw new Error(`No token in auth response. Keys: ${Object.keys(data).join(', ')}`);
  return token;
}

export function createAuthCurl(
  apiUrl: string,
  token: string,
): (method: string, path: string, opts?: { body?: string; headers?: Record<string, string>; timeout?: number }) => Promise<CurlResponse> {
  return (method, path, opts) =>
    curl(method, new URL(path, apiUrl).toString(), {
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
