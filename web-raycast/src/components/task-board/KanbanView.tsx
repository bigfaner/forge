import type { Task, TaskStatus } from '../../lib/types'
import { TaskCard } from '../TaskCard'
import { STATUS_LABEL } from '../../lib/utils'

const COLUMNS: TaskStatus[] = ['pending', 'in_progress', 'completed', 'blocked', 'skipped']

const COL_ACCENT: Record<TaskStatus, string> = {
  pending:     '#9c9c9d',
  in_progress: 'hsl(202,100%,67%)',
  completed:   'hsl(151,59%,59%)',
  blocked:     '#FF6363',
  skipped:     'hsl(43,100%,60%)',
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
    <div style={{ display: 'flex', gap: 12, height: '100%', minHeight: 0, overflowX: 'auto', paddingBottom: 8 }}>
      {activeCols.map(status => (
        <div key={status} style={{ display: 'flex', flexDirection: 'column', minWidth: 212, width: 212, flexShrink: 0 }}>
          {/* Column header */}
          <div style={{
            background: '#101111',
            border: '1px solid rgba(255,255,255,0.06)',
            borderBottom: 'none',
            borderRadius: '8px 8px 0 0',
            padding: '8px 12px',
            borderTop: `2px solid ${COL_ACCENT[status]}`,
          }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <span style={{ fontSize: 12, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px' }}>
                {STATUS_LABEL[status]}
              </span>
              <span style={{ fontSize: 11, color: '#434345', fontVariantNumeric: 'tabular-nums' }}>
                {grouped[status].length}
              </span>
            </div>
          </div>

          {/* Cards */}
          <div style={{
            flex: 1,
            background: 'rgba(255,255,255,0.02)',
            border: '1px solid rgba(255,255,255,0.06)',
            borderTop: 'none',
            borderRadius: '0 0 8px 8px',
            padding: 8,
            display: 'flex',
            flexDirection: 'column',
            gap: 6,
            overflowY: 'auto',
          }}>
            {grouped[status].map(t => (
              <TaskCard key={t.id} task={t} featureSlug={featureSlug} />
            ))}
            {grouped[status].length === 0 && (
              <p style={{ fontSize: 12, color: '#434345', textAlign: 'center', padding: '16px 0', letterSpacing: '0.2px' }}>
                Empty
              </p>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
