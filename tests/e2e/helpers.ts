import { execSync } from 'node:child_process';
import { readFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parse as parseYaml } from 'yaml';
import type { Page } from '@playwright/test';

const __dirname = dirname(fileURLToPath(import.meta.url));

// ── Config ─────────────────────────────────────────────────────────
// Lazy-loaded: only reads config.yaml when UI/API helpers are first called.
// CLI-only projects can omit config.yaml entirely.
let _configPath: string | null = null;
let _config: E2EConfig | null = null;

function findConfigPath(): string {
  // Allow explicit override via environment variable
  const envPath = process.env.E2E_CONFIG_PATH;
  if (envPath && existsSync(envPath)) return resolve(envPath);

  let dir = __dirname;
  for (let i = 0; i < 10; i++) {
    const candidate = resolve(dir, 'config.yaml');
    if (existsSync(candidate)) return candidate;
    const parent = resolve(dir, '..');
    if (parent === dir) break;
    dir = parent;
  }
  return ''; // CLI-only projects may not have config.yaml
}

// Screenshots go to <helpers-dir>/results/screenshots
const SCREENSHOTS_DIR = join(__dirname, 'results', 'screenshots');

interface E2EConfig {
  baseUrl?: string;
  apiBaseUrl?: string;
  timeout?: number | string;
  username?: string;
  password?: string;
  loginLocators?: { usernameField?: string; passwordField?: string; submitButton?: string };
}

function getConfig(): E2EConfig {
  if (_config) return _config;
  _configPath = findConfigPath();
  if (!_configPath) return {};
  const raw = parseYaml(readFileSync(_configPath, 'utf-8'));
  if (typeof raw !== 'object' || raw === null) {
    throw new Error(`Invalid config.yaml: expected object, got ${typeof raw}`);
  }
  _config = raw as E2EConfig;
  return _config;
}

export function baseUrl(): string { return getConfig().baseUrl ?? 'http://localhost:3456'; }
export function apiBaseUrl(): string { return getConfig().apiBaseUrl ?? 'http://localhost:8080'; }
const DEFAULT_TIMEOUT = 30000;

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

let _defaultCreds: UICredentials | null = null;
export function getDefaultCreds(): UICredentials {
  if (_defaultCreds) return _defaultCreds;
  _defaultCreds = {
    username: getConfig().username ?? 'admin',
    password: getConfig().password ?? 'password',
  };
  return _defaultCreds;
}
/** @deprecated Use getDefaultCreds() for lazy evaluation */
export const defaultCreds: UICredentials = new Proxy({} as UICredentials, {
  get(_, prop) { return getDefaultCreds()[prop as keyof UICredentials]; },
});

export async function loginViaUI(page: Page, creds: UICredentials = defaultCreds): Promise<void> {
  const loginUrl = new URL('/login', baseUrl()).toString();
  await page.goto(loginUrl);
  const locators = getConfig().loginLocators;
  const uPat = new RegExp(locators?.usernameField ?? 'username|email', 'i');
  const pPat = new RegExp(locators?.passwordField ?? 'password', 'i');
  const bPat = new RegExp(locators?.submitButton ?? 'login|sign in|submit', 'i');
  await page.getByRole('textbox', { name: uPat }).fill(creds.username);
  await page.getByRole('textbox', { name: pPat }).fill(creds.password);
  await page.getByRole('button', { name: bPat }).click();
  await page.waitForURL((url) => !url.pathname.includes('login') && url.pathname !== '/', { timeout: DEFAULT_TIMEOUT });
}

export async function getApiToken(apiBaseUrl: string, authPath: string, creds: UICredentials = defaultCreds): Promise<string> {
  const res = await curl('POST', `${apiBaseUrl}${authPath}`, {
    body: JSON.stringify({ username: creds.username, password: creds.password }),
  });
  if (res.status !== 200) throw new Error(`Auth failed: ${res.status} ${res.body}`);
  const data = JSON.parse(res.body);
  const token = data.token ?? data.access_token ?? data.data?.token;
  if (!token) throw new Error(`No token in auth response. Keys: ${Object.keys(data).join(', ')}`);
  return token;
}

export function createAuthCurl(
  apiBaseUrl: string,
  token: string,
): (method: string, path: string, opts?: { body?: string; headers?: Record<string, string>; timeout?: number }) => Promise<CurlResponse> {
  return (method, path, opts) => {
    const normalizedPath = path.startsWith('/') ? path : `/${path}`;
    return curl(method, `${apiBaseUrl}${normalizedPath}`, {
      ...opts,
      headers: { Authorization: `Bearer ${token}`, ...opts?.headers },
    });
  };
}

// ── CLI ────────────────────────────────────────────────────────────
export const PROJECT_ROOT = resolve(__dirname, '..', '..');

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
      cwd: cwd ?? PROJECT_ROOT,
    });
    return { stdout, stderr: '', exitCode: 0 };
  } catch (e: unknown) {
    const err = e as { stdout?: string; stderr?: string; status?: number };
    return {
      stdout: err.stdout ?? '',
      stderr: err.stderr ?? '',
      exitCode: err.status ?? 1,
    };
  }
}

// ── File helpers ───────────────────────────────────────────────────
export function readProjectFile(relPath: string): string {
  return readFileSync(join(PROJECT_ROOT, relPath), 'utf-8');
}

export function projectFileExists(relPath: string): boolean {
  return existsSync(join(PROJECT_ROOT, relPath));
}
