import { Link } from 'react-router-dom'
import type { Task } from '../lib/types'
import { StatusBadge } from './StatusBadge'

interface Props {
  task: Task
  featureSlug: string
  className?: string
}

export function TaskCard({ task, featureSlug }: Props) {
  return (
    <Link
      to={`/features/${featureSlug}/tasks/${task.id}`}
      style={{
        display: 'block',
        background: '#101111',
        border: '1px solid rgba(255,255,255,0.06)',
        borderRadius: 8,
        padding: '10px 12px',
        textDecoration: 'none',
        color: 'inherit',
        transition: 'border-color 0.15s',
      }}
      onMouseEnter={e => (e.currentTarget.style.borderColor = 'rgba(255,255,255,0.12)')}
      onMouseLeave={e => (e.currentTarget.style.borderColor = 'rgba(255,255,255,0.06)')}
    >
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 6, marginBottom: 6 }}>
        <span style={{ fontFamily: 'monospace', fontSize: 10, color: '#434345', letterSpacing: '0.3px' }}>
          {task.id}
        </span>
        <StatusBadge priority={task.priority} />
      </div>
      <p style={{ fontSize: 13, fontWeight: 500, color: '#f9f9f9', lineHeight: 1.4, letterSpacing: '0.2px', marginBottom: task.dependencies.length > 0 ? 6 : 0 }}>
        {task.title}
      </p>
      {task.dependencies.length > 0 && (
        <p style={{ fontSize: 11, color: '#434345', letterSpacing: '0.2px', fontFamily: 'monospace' }}>
          Deps: {task.dependencies.join(', ')}
        </p>
      )}
    </Link>
  )
}
