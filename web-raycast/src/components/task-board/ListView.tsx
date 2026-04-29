import { Link } from 'react-router-dom'
import type { Task } from '../../lib/types'
import { StatusBadge } from '../StatusBadge'

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
    <div style={{
      background: '#101111',
      border: '1px solid rgba(255,255,255,0.06)',
      borderRadius: 10,
      overflow: 'hidden',
    }}>
      <table style={{ width: '100%', fontSize: 13, borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ background: 'rgba(255,255,255,0.03)', borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              ID
            </th>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              Title
            </th>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              Phase
            </th>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              Priority
            </th>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              Status
            </th>
            <th style={{ textAlign: 'left', padding: '10px 12px', fontSize: 11, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px' }}>
              Deps
            </th>
          </tr>
        </thead>
        <tbody>
          {sorted.map((t, i) => (
            <tr
              key={t.id}
              style={{
                background: i % 2 === 0 ? 'rgba(255,255,255,0.02)' : 'rgba(255,255,255,0.005)',
                borderBottom: i < sorted.length - 1 ? '1px solid rgba(255,255,255,0.06)' : 'none',
                transition: 'background 0.15s',
              }}
              onMouseEnter={e => (e.currentTarget.style.background = 'rgba(255,255,255,0.05)')}
              onMouseLeave={e => (e.currentTarget.style.background = i % 2 === 0 ? 'rgba(255,255,255,0.02)' : 'rgba(255,255,255,0.005)')}
            >
              <td style={{ padding: '10px 12px' }}>
                <Link
                  to={`/features/${featureSlug}/tasks/${t.id}`}
                  style={{ fontSize: 11, color: '#434345', textDecoration: 'none', fontFamily: 'monospace', transition: 'color 0.15s' }}
                  onMouseEnter={e => (e.currentTarget.style.color = '#FF6363')}
                  onMouseLeave={e => (e.currentTarget.style.color = '#434345')}
                >
                  {t.id}
                </Link>
              </td>
              <td style={{ padding: '10px 12px' }}>
                <Link
                  to={`/features/${featureSlug}/tasks/${t.id}`}
                  style={{ color: '#f9f9f9', textDecoration: 'none', transition: 'color 0.15s', letterSpacing: '0.2px' }}
                  onMouseEnter={e => (e.currentTarget.style.color = '#FF6363')}
                  onMouseLeave={e => (e.currentTarget.style.color = '#f9f9f9')}
                >
                  {t.title}
                </Link>
              </td>
              <td style={{ padding: '10px 12px', color: '#6a6b6c', fontSize: 12 }}>{t.phase}</td>
              <td style={{ padding: '10px 12px' }}>
                <StatusBadge priority={t.priority} />
              </td>
              <td style={{ padding: '10px 12px' }}>
                <StatusBadge status={t.status} />
              </td>
              <td style={{ padding: '10px 12px', fontSize: 11, color: '#434345', fontFamily: 'monospace' }}>
                {t.dependencies.join(', ') || '—'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
