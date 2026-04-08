import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ChevronLeft } from 'lucide-react'
import { api } from '../lib/api'
import { MarkdownViewer } from '../components/MarkdownViewer'
import { KanbanView } from '../components/task-board/KanbanView'
import { ListView } from '../components/task-board/ListView'
import { DagView } from '../components/task-board/DagView'
import { ClaimButton } from '../components/ClaimButton'
import { cn } from '../lib/utils'

type Tab = 'prd' | 'design' | 'tasks'
type TaskView = 'kanban' | 'list' | 'dag'

export function FeatureDetail() {
  const { slug } = useParams<{ slug: string }>()
  const [tab, setTab] = useState<Tab>('tasks')
  const [taskView, setTaskView] = useState<TaskView>('kanban')

  const tasksQuery = useQuery({
    queryKey: ['features', slug, 'tasks'],
    queryFn: () => api.features.tasks(slug!),
    enabled: !!slug,
    refetchInterval: 30_000,
  })

  const prdQuery = useQuery({
    queryKey: ['features', slug, 'prd'],
    queryFn: () => api.features.prd(slug!),
    enabled: !!slug && tab === 'prd',
    staleTime: 60_000,
  })

  const designQuery = useQuery({
    queryKey: ['features', slug, 'design'],
    queryFn: () => api.features.design(slug!),
    enabled: !!slug && tab === 'design',
    staleTime: 60_000,
  })

  const tasks = tasksQuery.data?.tasks ?? []
  const TABS: { value: Tab; label: string }[] = [
    { value: 'tasks', label: 'Tasks' },
    { value: 'prd', label: 'PRD' },
    { value: 'design', label: 'Design' },
  ]
  const VIEWS: { value: TaskView; label: string }[] = [
    { value: 'kanban', label: 'Kanban' },
    { value: 'list', label: 'List' },
    { value: 'dag', label: 'DAG' },
  ]

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border px-6 py-4">
        <Link to="/features" className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground mb-2 w-fit">
          <ChevronLeft className="h-3 w-3" /> Features
        </Link>
        <div className="flex items-center justify-between">
          <div>
            <p className="font-mono text-xs text-muted-foreground">{slug}</p>
            <h1 className="text-lg font-semibold">{slug}</h1>
          </div>
          {tab === 'tasks' && <ClaimButton featureSlug={slug!} />}
        </div>

        {/* Tabs */}
        <div className="flex gap-0 mt-4 border-b border-transparent -mb-px">
          {TABS.map(t => (
            <button
              key={t.value}
              onClick={() => setTab(t.value)}
              className={cn(
                'px-4 py-2 text-sm font-medium border-b-2 transition-colors cursor-pointer',
                tab === t.value
                  ? 'border-accent text-accent'
                  : 'border-transparent text-muted-foreground hover:text-foreground',
              )}
            >
              {t.label}
              {t.value === 'tasks' && tasks.length > 0 && (
                <span className="ml-1.5 text-xs text-muted-foreground">({tasks.length})</span>
              )}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {tab === 'tasks' && (
          <div className="h-full flex flex-col">
            {/* View switcher */}
            <div className="flex gap-1 p-3 border-b border-border bg-background/50">
              {VIEWS.map(v => (
                <button
                  key={v.value}
                  onClick={() => setTaskView(v.value)}
                  className={cn(
                    'px-3 py-1 rounded text-xs font-medium transition-colors cursor-pointer',
                    taskView === v.value
                      ? 'bg-accent/10 text-accent'
                      : 'text-muted-foreground hover:text-foreground',
                  )}
                >
                  {v.label}
                </button>
              ))}
            </div>
            <div className="flex-1 overflow-auto p-4">
              {tasksQuery.isLoading ? (
                <div className="flex items-center justify-center h-32 text-sm text-muted-foreground">Loading tasks…</div>
              ) : taskView === 'kanban' ? (
                <KanbanView tasks={tasks} featureSlug={slug!} />
              ) : taskView === 'list' ? (
                <ListView tasks={tasks} featureSlug={slug!} />
              ) : (
                <DagView tasks={tasks} featureSlug={slug!} />
              )}
            </div>
          </div>
        )}

        {tab === 'prd' && (
          <div className="p-6">
            {prdQuery.isLoading ? (
              <div className="text-sm text-muted-foreground">Loading PRD…</div>
            ) : prdQuery.data ? (
              <MarkdownViewer content={prdQuery.data} />
            ) : (
              <p className="text-sm text-muted-foreground">PRD not found for this feature.</p>
            )}
          </div>
        )}

        {tab === 'design' && (
          <div className="p-6">
            {designQuery.isLoading ? (
              <div className="text-sm text-muted-foreground">Loading Design…</div>
            ) : designQuery.data ? (
              <MarkdownViewer content={designQuery.data} />
            ) : (
              <p className="text-sm text-muted-foreground">Design doc not found for this feature.</p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
