import { useQuery } from '@tanstack/react-query'
import { GitCommit } from 'lucide-react'
import { api } from '../lib/api'
import { formatDate } from '../lib/utils'

export function Records() {
  const { data, isLoading } = useQuery({
    queryKey: ['records'],
    queryFn: () => api.records.all(),
    refetchInterval: 30_000,
  })

  const records = data?.records ?? []

  return (
    <div className="p-6 space-y-5">
      <div>
        <h1 className="text-xl font-semibold">Execution Records</h1>
        <p className="text-sm text-muted-foreground mt-0.5">Timeline of all completed tasks</p>
      </div>

      {isLoading ? (
        <div className="space-y-2">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="h-16 rounded-lg bg-muted animate-pulse" />
          ))}
        </div>
      ) : records.length === 0 ? (
        <div className="rounded-lg border border-dashed border-border p-12 text-center">
          <p className="text-sm text-muted-foreground">No execution records yet</p>
        </div>
      ) : (
        <div className="relative pl-4">
          {/* Timeline line */}
          <div className="absolute left-0 top-0 bottom-0 w-px bg-border ml-1.5" />

          <div className="space-y-3">
            {records.map((r, i) => (
              <div key={i} className="relative pl-5">
                {/* Dot */}
                <div className="absolute left-0 top-2 h-2 w-2 rounded-full bg-accent ring-2 ring-background" />

                <div className="rounded-lg border border-border bg-card p-3">
                  <div className="flex items-start justify-between gap-2 mb-1">
                    <div>
                      <span className="font-mono text-xs text-muted-foreground">{r.featureSlug} / {r.taskId}</span>
                      <p className="text-sm font-medium">{r.taskTitle}</p>
                    </div>
                    <span className="text-xs text-muted-foreground shrink-0">{formatDate(r.completedAt)}</span>
                  </div>
                  <div className="flex items-center gap-3 text-xs text-muted-foreground">
                    {r.coverage && <span className="text-green-400">Coverage {r.coverage}</span>}
                    {r.filesChanged > 0 && <span>{r.filesChanged} files changed</span>}
                    {r.commitHash && (
                      <span className="flex items-center gap-1 font-mono">
                        <GitCommit className="h-3 w-3" />
                        {r.commitHash.slice(0, 7)}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
