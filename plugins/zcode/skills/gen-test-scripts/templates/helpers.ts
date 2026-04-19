import { execSync } from 'node:child_process';
import { writeFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const DEFAULT_TIMEOUT = 30000;
const SCREENSHOTS_DIR = join(__dirname, '..', 'results', 'screenshots');

export const baseUrl = process.env.E2E_BASE_URL ?? 'http://localhost:3456';

export function ab(cmd: string): string {
  return execSync(`agent-browser ${cmd}`, {
    encoding: 'utf-8',
    timeout: DEFAULT_TIMEOUT,
  });
}

export function abJson(cmd: string): any {
  const raw = ab(`${cmd} --json`);
  try {
    return JSON.parse(raw);
  } catch {
    throw new Error(`Failed to parse agent-browser JSON output: ${raw.slice(0, 200)}`);
  }
}

export function snapshotContains(text: string): boolean {
  const result = abJson('snapshot');
  return result?.data?.snapshot?.includes(text) ?? false;
}

export function findElement(role: string, name?: string): string | null {
  const cmd = name
    ? `find role ${role} --name "${name}" --json`
    : `find role ${role} --json`;
  try {
    const result = abJson(cmd);
    return result?.data?.ref ?? result?.ref ?? null;
  } catch {
    return null;
  }
}

export function screenshot(tcId: string): string {
  if (!existsSync(SCREENSHOTS_DIR)) {
    mkdirSync(SCREENSHOTS_DIR, { recursive: true });
  }
  const path = join(SCREENSHOTS_DIR, `${tcId}.png`);
  ab(`screenshot "${path}"`);
  return path;
}

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
