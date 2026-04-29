import type { TaskStatus, TaskPriority } from '../lib/types'
import { STATUS_LABEL } from '../lib/utils'

const STATUS_STYLE: Record<TaskStatus, { bg: string; color: string; border: string }> = {
  pending:     { bg: 'rgba(156,156,157,0.1)', color: '#9c9c9d', border: 'rgba(156,156,157,0.25)' },
  in_progress: { bg: 'rgba(85,179,255,0.12)', color: 'hsl(202,100%,67%)', border: 'rgba(85,179,255,0.3)' },
  completed:   { bg: 'rgba(95,201,146,0.12)', color: 'hsl(151,59%,59%)', border: 'rgba(95,201,146,0.3)' },
  blocked:     { bg: 'rgba(255,99,99,0.12)',  color: '#FF6363', border: 'rgba(255,99,99,0.3)' },
  skipped:     { bg: 'rgba(255,188,51,0.12)', color: 'hsl(43,100%,60%)', border: 'rgba(255,188,51,0.3)' },
}

const PRIORITY_STYLE: Record<TaskPriority, { bg: string; color: string; border: string }> = {
  P0: { bg: 'rgba(255,99,99,0.12)',  color: '#FF6363', border: 'rgba(255,99,99,0.3)' },
  P1: { bg: 'rgba(255,165,0,0.12)',  color: '#ffa500', border: 'rgba(255,165,0,0.3)' },
  P2: { bg: 'rgba(156,156,157,0.1)', color: '#6a6b6c', border: 'rgba(156,156,157,0.2)' },
}

interface Props {
  status?: TaskStatus
  priority?: TaskPriority
  className?: string
}

export function StatusBadge({ status, priority }: Props) {
  if (status) {
    const s = STATUS_STYLE[status]
    return (
      <span style={{
        display: 'inline-flex', alignItems: 'center',
        padding: '2px 7px', borderRadius: 4,
        fontSize: 11, fontWeight: 600, letterSpacing: '0.2px',
        background: s.bg, color: s.color,
        border: `1px solid ${s.border}`,
      }}>
        {STATUS_LABEL[status]}
      </span>
    )
  }
  if (priority) {
    const p = PRIORITY_STYLE[priority]
    return (
      <span style={{
        display: 'inline-flex', alignItems: 'center',
        padding: '2px 7px', borderRadius: 4,
        fontSize: 11, fontWeight: 600, letterSpacing: '0.2px',
        background: p.bg, color: p.color,
        border: `1px solid ${p.border}`,
      }}>
        {priority}
      </span>
    )
  }
  return null
}
