import type { TaskStatus, TaskPriority } from '../lib/types'
import { STATUS_COLOR, STATUS_LABEL, PRIORITY_COLOR, cn } from '../lib/utils'

interface Props {
  status?: TaskStatus
  priority?: TaskPriority
  className?: string
}

export function StatusBadge({ status, priority, className }: Props) {
  if (status) {
    return (
      <span className={cn('inline-flex items-center rounded border px-1.5 py-0.5 text-xs font-medium', STATUS_COLOR[status], className)}>
        {STATUS_LABEL[status]}
      </span>
    )
  }
  if (priority) {
    return (
      <span className={cn('inline-flex items-center rounded border px-1.5 py-0.5 text-xs font-medium', PRIORITY_COLOR[priority], className)}>
        {priority}
      </span>
    )
  }
  return null
}
