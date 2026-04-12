import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ChevronLeft, FileCode, GitCommit } from 'lucide-react'
import { api } from '../lib/api'
import { StatusBadge } from '../components/StatusBadge'
import { formatDate } from '../lib/utils'
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

  if (isLoading) return <div style={{ padding: 32, fontSize: 14, color: '#6a6b6c' }}>Loading task…</div>
  if (!task) return <div style={{ padding: 32, fontSize: 14, color: '#6a6b6c' }}>Task not found.</div>

  const sectionLabel = (text: string) => (
    <p style={{ fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.6px', textTransform: 'uppercase', marginBottom: 8 }}>
      {text}
    </p>
  )

  return (
    <div style={{ padding: '32px 40px', maxWidth: 780 }}>
      <Link
        to={`/features/${slug}`}
        style={{ display: 'inline-flex', alignItems: 'center', gap: 4, fontSize: 12, color: '#6a6b6c', textDecoration: 'none', letterSpacing: '0.2px', marginBottom: 20, transition: 'color 0.15s' }}
        onMouseEnter={e => (e.currentTarget.style.color = '#f9f9f9')}
        onMouseLeave={e => (e.currentTarget.style.color = '#6a6b6c')}
      >
        <ChevronLeft size={12} /> {slug}
      </Link>

      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 16, marginBottom: 28 }}>
        <div>
          <p style={{ fontFamily: 'monospace', fontSize: 11, color: '#434345', letterSpacing: '0.3px', marginBottom: 6 }}>{task.id}</p>
          <h1 style={{ fontSize: 20, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0, lineHeight: 1.3 }}>{task.title}</h1>
        </div>
        <div style={{ display: 'flex', gap: 6, flexShrink: 0 }}>
          <StatusBadge status={task.status} />
          <StatusBadge priority={task.priority} />
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 8, marginBottom: 28 }}>
        {[
          { label: 'Phase', value: `Phase ${task.phase}` },
          { label: 'Estimated', value: task.estimatedTime || '—' },
          { label: 'Dependencies', value: task.dependencies.length > 0 ? task.dependencies.join(', ') : 'none' },
          { label: 'Files', value: `${task.files.length} file(s)` },
        ].map(({ label, value }) => (
          <div key={label} style={{ background: '#101111', border: '1px solid rgba(255,255,255,0.06)', borderRadius: 8, padding: '10px 14px' }}>
            <p style={{ fontSize: 11, color: '#6a6b6c', letterSpacing: '0.2px', marginBottom: 4 }}>{label}</p>
            <p style={{ fontSize: 13, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>{value}</p>
          </div>
        ))}
      </div>

      {task.description && (
        <p style={{ fontSize: 14, color: '#9c9c9d', lineHeight: 1.65, letterSpacing: '0.2px', marginBottom: 28 }}>
          {task.description}
        </p>
      )}

      {task.files.length > 0 && (
        <div style={{ marginBottom: 28 }}>
          {sectionLabel('Files')}
          <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            {task.files.map(f => (
              <div key={f} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 12, color: '#6a6b6c', fontFamily: 'monospace' }}>
                <FileCode size={12} />
                {f}
              </div>
            ))}
          </div>
        </div>
      )}

      {task.status !== 'completed' && (
        <div style={{ marginBottom: 28 }}>
          {sectionLabel('Update Status')}
          <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
            {MUTABLE_STATUSES.filter(s => s !== task.status).map(s => (
              <button
                key={s}
                onClick={() => setStatus.mutate(s)}
                disabled={setStatus.isPending}
                style={{
                  padding: '6px 14px', borderRadius: 6, fontSize: 13, fontWeight: 500,
                  letterSpacing: '0.2px', cursor: 'pointer', border: '1px solid rgba(255,255,255,0.08)',
                  background: 'rgba(255,255,255,0.04)', color: '#9c9c9d', transition: 'all 0.15s',
                  fontFamily: 'inherit', opacity: setStatus.isPending ? 0.5 : 1,
                }}
                onMouseEnter={e => { e.currentTarget.style.color = '#f9f9f9'; e.currentTarget.style.borderColor = 'rgba(255,255,255,0.16)' }}
                onMouseLeave={e => { e.currentTarget.style.color = '#9c9c9d'; e.currentTarget.style.borderColor = 'rgba(255,255,255,0.08)' }}
              >
                → {s.replace('_', ' ')}
              </button>
            ))}
          </div>
        </div>
      )}

      {task.record && typeof task.record === 'object' && (
        <div style={{ background: '#101111', border: '1px solid rgba(255,255,255,0.06)', boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset', borderRadius: 12, padding: '20px 24px' }}>
          {sectionLabel('Execution Record')}
          <p style={{ fontSize: 14, color: '#cecece', lineHeight: 1.65, letterSpacing: '0.2px', marginBottom: 16 }}>
            {task.record.summary}
          </p>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8, marginBottom: 16 }}>
            {[
              { label: 'Coverage', value: task.record.coverage || '—' },
              { label: 'Tests', value: task.record.testResults || '—' },
              { label: 'Completed', value: formatDate(task.record.completedAt) },
            ].map(({ label, value }) => (
              <div key={label} style={{ background: 'rgba(255,255,255,0.03)', borderRadius: 6, padding: '8px 12px' }}>
                <p style={{ fontSize: 11, color: '#6a6b6c', letterSpacing: '0.2px', marginBottom: 3 }}>{label}</p>
                <p style={{ fontSize: 13, fontWeight: 500, color: '#f9f9f9', margin: 0 }}>{value}</p>
              </div>
            ))}
          </div>
          {task.record.decisions?.length > 0 && (
            <div style={{ marginBottom: 12 }}>
              <p style={{ fontSize: 11, color: '#6a6b6c', letterSpacing: '0.2px', marginBottom: 6 }}>Decisions</p>
              <ul style={{ margin: 0, padding: 0, listStyle: 'none', display: 'flex', flexDirection: 'column', gap: 3 }}>
                {task.record.decisions.map((d, i) => (
                  <li key={i} style={{ fontSize: 12, color: '#9c9c9d', letterSpacing: '0.2px' }}>• {d}</li>
                ))}
              </ul>
            </div>
          )}
          {task.record.commitHash && (
            <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 11, color: '#434345', fontFamily: 'monospace' }}>
              <GitCommit size={11} />
              {task.record.commitHash}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
