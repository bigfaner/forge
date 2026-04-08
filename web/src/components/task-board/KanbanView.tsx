import type { Task, TaskStatus } from '../../lib/types'
import { TaskCard } from '../TaskCard'
import { STATUS_LABEL } from '../../lib/utils'
import { cn } from '../../lib/utils'

const COLUMNS: TaskStatus[] = ['pending', 'in_progress', 'completed', 'blocked', 'skipped']

const COLUMN_STYLE: Record<TaskStatus, string> = {
  pending:     'border-t-slate-500',
  in_progress: 'border-t-blue-500',
  completed:   'border-t-green-500',
  blocked:     'border-t-red-500',
  skipped:     'border-t-yellow-500',
}

interface Props {
  tasks: Task[]
  featureSlug: string
}

export function KanbanView({ tasks, featureSlug }: Props) {
  const grouped = COLUMNS.reduce<Record<TaskStatus, Task[]>>(
    (acc, s) => ({ ...acc, [s]: [] }),
    {} as Record<TaskStatus, Task[]>,
  )
  tasks.forEach(t => grouped[t.status]?.push(t))

  const activeCols = COLUMNS.filter(s => grouped[s].length > 0 || s === 'pending' || s === 'in_progress')

  return (
    <div className="flex gap-3 h-full min-h-0 overflow-x-auto pb-2">
      {activeCols.map(status => (
        <div key={status} className={cn('flex flex-col min-w-52 w-52 shrink-0')}>
          <div className={cn('rounded-t-lg border-t-2 border-x border-border bg-card px-3 py-2', COLUMN_STYLE[status])}>
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium">{STATUS_LABEL[status]}</span>
              <span className="text-xs text-muted-foreground">{grouped[status].length}</span>
            </div>
          </div>
          <div className="flex-1 rounded-b-lg border border-t-0 border-border bg-card/50 p-2 space-y-2 overflow-y-auto">
            {grouped[status].map(t => (
              <TaskCard key={t.id} task={t} featureSlug={featureSlug} />
            ))}
            {grouped[status].length === 0 && (
              <p className="text-xs text-muted-foreground text-center py-4">Empty</p>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
