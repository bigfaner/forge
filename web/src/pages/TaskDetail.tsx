import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ChevronLeft, FileCode, GitCommit } from 'lucide-react'
import { api } from '../lib/api'
import { StatusBadge } from '../components/StatusBadge'
import { formatDate, cn } from '../lib/utils'
import type { TaskStatus } from '../lib/types'

const MUTABLE_STATUSES: TaskStatus[] = ['blocked', 'skipped', 'pending', 'in_progress']

export function TaskDetail() {
  const { slug, id } = useParams<{ slug: string; id: string }>()
  const qc = useQueryClient()

  const { data: task, isLoading } = useQuery({
    queryKey: ['features', slug, 'tasks', id],
    queryFn: () => api.features.task(slug!, id!),
    enabled: !!slug && !!id,
  })

  const setStatus = useMutation({
    mutationFn: (status: TaskStatus) => api.tasks.setStatus(slug!, id!, status),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['features', slug, 'tasks'] })
      qc.invalidateQueries({ queryKey: ['features', slug, 'tasks', id] })
    },
  })

  if (isLoading) {
    return <div className="p-6 text-sm text-muted-foreground">Loading task…</div>
  }
  if (!task) {
    return <div className="p-6 text-sm text-muted-foreground">Task not found.</div>
  }

  return (
    <div className="p-6 max-w-3xl space-y-6">
      <div>
        <Link to={`/features/${slug}`} className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground mb-3 w-fit">
          <ChevronLeft className="h-3 w-3" /> {slug}
        </Link>
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="font-mono text-xs text-muted-foreground mb-1">{task.id}</p>
            <h1 className="text-lg font-semibold">{task.title}</h1>
          </div>
          <div className="flex gap-2 shrink-0">
            <StatusBadge status={task.status} />
            <StatusBadge priority={task.priority} />
          </div>
        </div>
      </div>

      {/* Meta */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 text-sm">
        <MetaField label="Phase" value={`Phase ${task.phase}`} />
        <MetaField label="Estimated" value={task.estimatedTime || '—'} />
        <MetaField label="Dependencies" value={task.dependencies.length > 0 ? task.dependencies.join(', ') : 'none'} />
        <MetaField label="Files" value={`${task.files.length} file(s)`} />
      </div>

      {task.description && (
        <p className="text-sm text-muted-foreground leading-relaxed">{task.description}</p>
      )}

      {/* File list */}
      {task.files.length > 0 && (
        <div>
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">Files</p>
          <div className="space-y-1">
            {task.files.map(f => (
              <div key={f} className="flex items-center gap-2 text-xs text-muted-foreground font-mono">
                <FileCode className="h-3 w-3 shrink-0" />
                {f}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Status update */}
      {task.status !== 'completed' && (
        <div>
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">Update Status</p>
          <div className="flex gap-2 flex-wrap">
            {MUTABLE_STATUSES.filter(s => s !== task.status).map(s => (
              <button
                key={s}
                onClick={() => setStatus.mutate(s)}
                disabled={setStatus.isPending}
                className={cn(
                  'px-3 py-1.5 rounded-md border text-xs font-medium transition-colors cursor-pointer',
                  'border-border bg-card hover:bg-muted disabled:opacity-50',
                )}
              >
                → {s.replace('_', ' ')}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Execution record */}
      {task.record && typeof task.record === 'object' && (
        <div className="rounded-lg border border-border bg-card p-4 space-y-4">
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Execution Record</p>

          <p className="text-sm">{task.record.summary}</p>

          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 text-sm">
            <MetaField label="Coverage" value={task.record.coverage || '—'} />
            <MetaField label="Tests" value={task.record.testResults || '—'} />
            <MetaField label="Completed" value={formatDate(task.record.completedAt)} />
          </div>

          {task.record.decisions?.length > 0 && (
            <div>
              <p className="text-xs text-muted-foreground mb-1">Decisions</p>
              <ul className="space-y-0.5">
                {task.record.decisions.map((d, i) => (
                  <li key={i} className="text-xs text-muted-foreground">• {d}</li>
                ))}
              </ul>
            </div>
          )}

          {task.record.commitHash && (
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-mono">
              <GitCommit className="h-3 w-3" />
              {task.record.commitHash}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function MetaField({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md bg-muted px-3 py-2">
      <p className="text-xs text-muted-foreground mb-0.5">{label}</p>
      <p className="text-sm font-medium">{value}</p>
    </div>
  )
}
