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
    <div style={{ padding: '32px 40px', maxWidth: 900 }}>
      <div style={{ marginBottom: 28 }}>
        <h1 style={{ fontSize: 24, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
          Execution Records
        </h1>
        <p style={{ fontSize: 14, color: '#6a6b6c', marginTop: 4, letterSpacing: '0.2px' }}>
          Timeline of all completed tasks
        </p>
      </div>

      {isLoading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {[...Array(5)].map((_, i) => (
            <div key={i} style={{ height: 64, borderRadius: 12, background: 'rgba(255,255,255,0.03)' }} />
          ))}
        </div>
      ) : records.length === 0 ? (
        <div style={{ border: '1px dashed rgba(255,255,255,0.1)', borderRadius: 12, padding: '64px 32px', textAlign: 'center' }}>
          <p style={{ fontSize: 14, color: '#6a6b6c', letterSpacing: '0.2px' }}>No execution records yet</p>
        </div>
      ) : (
        <div style={{ position: 'relative', paddingLeft: 20 }}>
          {/* Timeline line */}
          <div style={{
            position: 'absolute', left: 5, top: 8, bottom: 8,
            width: 1, background: 'rgba(255,255,255,0.06)',
          }} />

          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {records.map((r, i) => (
              <div key={i} style={{ position: 'relative', paddingLeft: 20 }}>
                {/* Timeline dot */}
                <div style={{
                  position: 'absolute', left: -1, top: 18,
                  height: 7, width: 7, borderRadius: '50%',
                  background: '#FF6363',
                  boxShadow: 'rgba(255,99,99,0.4) 0px 0px 6px',
                }} />

                <div style={{
                  background: '#101111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset',
                  borderRadius: 10,
                  padding: '12px 16px',
                }}>
                  <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 8, marginBottom: 6 }}>
                    <div>
                      <span style={{ fontFamily: 'monospace', fontSize: 11, color: '#434345', letterSpacing: '0.3px' }}>
                        {r.featureSlug} / {r.taskId}
                      </span>
                      <p style={{ fontSize: 13, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px', margin: '2px 0 0' }}>
                        {r.taskTitle}
                      </p>
                    </div>
                    <span style={{ fontSize: 11, color: '#434345', flexShrink: 0, letterSpacing: '0.2px' }}>
                      {formatDate(r.completedAt)}
                    </span>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                    {r.coverage && (
                      <span style={{ fontSize: 12, color: 'hsl(151,59%,59%)', letterSpacing: '0.2px' }}>
                        Coverage {r.coverage}
                      </span>
                    )}
                    {r.filesChanged > 0 && (
                      <span style={{ fontSize: 12, color: '#6a6b6c', letterSpacing: '0.2px' }}>
                        {r.filesChanged} files changed
                      </span>
                    )}
                    {r.commitHash && (
                      <span style={{ display: 'flex', alignItems: 'center', gap: 4, fontSize: 11, color: '#434345', fontFamily: 'monospace' }}>
                        <GitCommit size={11} />
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
