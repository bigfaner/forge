import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { TaskStatus, TaskPriority, LessonCategory } from './types'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const STATUS_LABEL: Record<TaskStatus, string> = {
  pending: 'Pending',
  in_progress: 'In Progress',
  completed: 'Completed',
  blocked: 'Blocked',
  skipped: 'Skipped',
}

export const STATUS_COLOR: Record<TaskStatus, string> = {
  pending:     'bg-slate-500/20 text-slate-400 border-slate-500/30',
  in_progress: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
  completed:   'bg-green-500/20 text-green-400 border-green-500/30',
  blocked:     'bg-red-500/20 text-red-400 border-red-500/30',
  skipped:     'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
}

export const PRIORITY_COLOR: Record<TaskPriority, string> = {
  P0: 'bg-red-500/20 text-red-400 border-red-500/30',
  P1: 'bg-orange-500/20 text-orange-400 border-orange-500/30',
  P2: 'bg-slate-500/20 text-slate-400 border-slate-500/30',
}

export const LESSON_CATEGORY_COLOR: Record<LessonCategory, string> = {
  debug:   'bg-red-500/20 text-red-400',
  arch:    'bg-purple-500/20 text-purple-400',
  tool:    'bg-blue-500/20 text-blue-400',
  pattern: 'bg-cyan-500/20 text-cyan-400',
  gotcha:  'bg-yellow-500/20 text-yellow-400',
  other:   'bg-slate-500/20 text-slate-400',
}

export function formatDate(iso: string): string {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export function progressPercent(stats: { completed: number; total: number }): number {
  if (stats.total === 0) return 0
  return Math.round((stats.completed / stats.total) * 100)
}
