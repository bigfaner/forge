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

type Tab = 'prd' | 'design' | 'tasks'
type TaskView = 'kanban' | 'list' | 'dag'

const TAB_TABS: { value: Tab; label: string }[] = [
  { value: 'tasks', label: 'Tasks' },
  { value: 'prd', label: 'PRD' },
  { value: 'design', label: 'Design' },
]

const VIEW_OPTS: { value: TaskView; label: string }[] = [
  { value: 'kanban', label: 'Kanban' },
  { value: 'list', label: 'List' },
  { value: 'dag', label: 'DAG' },
]

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

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      {/* Header */}
      <div style={{
        padding: '16px 32px 0',
        borderBottom: '1px solid rgba(255,255,255,0.06)',
        background: '#07080a',
      }}>
        <Link
          to="/features"
          style={{
            display: 'inline-flex', alignItems: 'center', gap: 4,
            fontSize: 12, color: '#6a6b6c', textDecoration: 'none',
            letterSpacing: '0.2px', marginBottom: 10,
            transition: 'color 0.15s',
          }}
          onMouseEnter={e => (e.currentTarget.style.color = '#f9f9f9')}
          onMouseLeave={e => (e.currentTarget.style.color = '#6a6b6c')}
        >
          <ChevronLeft size={12} /> Features
        </Link>

        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
          <div>
            <p style={{ fontFamily: 'monospace', fontSize: 11, color: '#434345', letterSpacing: '0.3px', marginBottom: 3 }}>
              {slug}
            </p>
            <h1 style={{ fontSize: 18, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
              {slug}
            </h1>
          </div>
          {tab === 'tasks' && <ClaimButton featureSlug={slug!} />}
        </div>

        {/* Tabs */}
        <div style={{ display: 'flex', gap: 0 }}>
          {TAB_TABS.map(t => (
            <button
              key={t.value}
              onClick={() => setTab(t.value)}
              style={{
                padding: '8px 16px',
                fontSize: 14,
                fontWeight: 500,
                letterSpacing: '0.2px',
                cursor: 'pointer',
                border: 'none',
                background: 'transparent',
                color: tab === t.value ? '#f9f9f9' : '#6a6b6c',
                borderBottom: `2px solid ${tab === t.value ? '#FF6363' : 'transparent'}`,
                transition: 'color 0.15s, border-color 0.15s',
                fontFamily: 'inherit',
              }}
            >
              {t.label}
              {t.value === 'tasks' && tasks.length > 0 && (
                <span style={{ marginLeft: 6, fontSize: 11, color: '#434345' }}>({tasks.length})</span>
              )}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div style={{ flex: 1, overflow: 'auto' }}>
        {tab === 'tasks' && (
          <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* View switcher */}
            <div style={{
              display: 'flex', gap: 4, padding: '10px 16px',
              borderBottom: '1px solid rgba(255,255,255,0.06)',
              background: 'rgba(255,255,255,0.02)',
            }}>
              {VIEW_OPTS.map(v => (
                <button
                  key={v.value}
                  onClick={() => setTaskView(v.value)}
                  style={{
                    padding: '4px 12px',
                    borderRadius: 6,
                    fontSize: 12,
                    fontWeight: 500,
                    letterSpacing: '0.2px',
                    cursor: 'pointer',
                    border: 'none',
                    background: taskView === v.value ? 'rgba(255,255,255,0.08)' : 'transparent',
                    color: taskView === v.value ? '#f9f9f9' : '#6a6b6c',
                    transition: 'all 0.15s',
                    fontFamily: 'inherit',
                  }}
                >
                  {v.label}
                </button>
              ))}
            </div>
            <div style={{ flex: 1, overflow: 'auto', padding: 16 }}>
              {tasksQuery.isLoading ? (
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: 128, fontSize: 14, color: '#6a6b6c' }}>
                  Loading tasks…
                </div>
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
          <div style={{ padding: 32 }}>
            {prdQuery.isLoading ? (
              <div style={{ fontSize: 14, color: '#6a6b6c' }}>Loading PRD…</div>
            ) : prdQuery.data ? (
              <MarkdownViewer content={prdQuery.data} />
            ) : (
              <p style={{ fontSize: 14, color: '#6a6b6c' }}>PRD not found for this feature.</p>
            )}
          </div>
        )}

        {tab === 'design' && (
          <div style={{ padding: 32 }}>
            {designQuery.isLoading ? (
              <div style={{ fontSize: 14, color: '#6a6b6c' }}>Loading Design…</div>
            ) : designQuery.data ? (
              <MarkdownViewer content={designQuery.data} />
            ) : (
              <p style={{ fontSize: 14, color: '#6a6b6c' }}>Design doc not found for this feature.</p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
