import { Link } from 'react-router-dom'
import type { Task } from '../../lib/types'
import { StatusBadge } from '../StatusBadge'
import { cn } from '../../lib/utils'

interface Props {
  tasks: Task[]
  featureSlug: string
}

export function ListView({ tasks, featureSlug }: Props) {
  const sorted = [...tasks].sort((a, b) => {
    const phaseD = a.phase - b.phase
    if (phaseD !== 0) return phaseD
    const pOrder = { P0: 0, P1: 1, P2: 2 }
    return pOrder[a.priority] - pOrder[b.priority]
  })

  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-border bg-muted/50">
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">ID</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">Title</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">Phase</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">Priority</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">Status</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-muted-foreground">Deps</th>
          </tr>
        </thead>
        <tbody>
          {sorted.map((t, i) => (
            <tr
              key={t.id}
              className={cn(
                'border-b border-border last:border-0 hover:bg-muted/30 transition-colors',
                i % 2 === 0 ? 'bg-card' : 'bg-background',
              )}
            >
              <td className="px-4 py-2.5">
                <Link to={`/features/${featureSlug}/tasks/${t.id}`} className="font-mono text-xs text-muted-foreground hover:text-accent">
                  {t.id}
                </Link>
              </td>
              <td className="px-4 py-2.5">
                <Link to={`/features/${featureSlug}/tasks/${t.id}`} className="hover:text-accent transition-colors">
                  {t.title}
                </Link>
              </td>
              <td className="px-4 py-2.5 text-muted-foreground text-xs">{t.phase}</td>
              <td className="px-4 py-2.5"><StatusBadge priority={t.priority} /></td>
              <td className="px-4 py-2.5"><StatusBadge status={t.status} /></td>
              <td className="px-4 py-2.5 text-xs text-muted-foreground font-mono">
                {t.dependencies.join(', ') || '—'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
