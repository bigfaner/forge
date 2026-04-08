import { Link } from 'react-router-dom'
import type { Task } from '../lib/types'
import { StatusBadge } from './StatusBadge'
import { cn } from '../lib/utils'

interface Props {
  task: Task
  featureSlug: string
  className?: string
}

export function TaskCard({ task, featureSlug, className }: Props) {
  return (
    <Link
      to={`/features/${featureSlug}/tasks/${task.id}`}
      className={cn(
        'block rounded-md border border-border bg-card p-3',
        'hover:border-accent/30 transition-colors',
        className,
      )}
    >
      <div className="flex items-start justify-between gap-2 mb-2">
        <span className="font-mono text-xs text-muted-foreground">{task.id}</span>
        <StatusBadge priority={task.priority} />
      </div>
      <p className="text-sm font-medium leading-snug mb-2">{task.title}</p>
      {task.dependencies.length > 0 && (
        <p className="text-xs text-muted-foreground">
          Deps: {task.dependencies.join(', ')}
        </p>
      )}
    </Link>
  )
}
