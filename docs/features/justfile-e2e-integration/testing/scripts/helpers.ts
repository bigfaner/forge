import { execSync } from 'node:child_process';
import { readFileSync, mkdirSync, existsSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

// Project root: docs/features/<slug>/testing/scripts/ → 5 levels up
export const PROJECT_ROOT = resolve(__dirname, '..', '..', '..', '..', '..');

const DEFAULT_TIMEOUT = 30000;

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
      cwd: cwd ?? PROJECT_ROOT,
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

// ── File helpers ───────────────────────────────────────────────────
export function readProjectFile(relPath: string): string {
  return readFileSync(join(PROJECT_ROOT, relPath), 'utf-8');
}

export function projectFileExists(relPath: string): boolean {
  return existsSync(join(PROJECT_ROOT, relPath));
}
