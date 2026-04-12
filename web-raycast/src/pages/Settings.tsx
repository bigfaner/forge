import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'

export function Settings() {
  const { data } = useQuery({
    queryKey: ['health'],
    queryFn: () => api.health(),
    refetchInterval: 60_000,
  })

  return (
    <div style={{ padding: '32px 40px', maxWidth: 560 }}>
      <div style={{ marginBottom: 28 }}>
        <h1 style={{ fontSize: 24, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
          Settings
        </h1>
        <p style={{ fontSize: 14, color: '#6a6b6c', marginTop: 4, letterSpacing: '0.2px' }}>
          Server info and configuration
        </p>
      </div>

      <div style={{
        background: '#101111',
        border: '1px solid rgba(255,255,255,0.06)',
        boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset',
        borderRadius: 12,
        overflow: 'hidden',
      }}>
        <InfoRow label="task-cli version" value={data?.version ?? '—'} />
        <InfoRow label="Project root" value={data?.projectRoot ?? '—'} mono />
        <InfoRow label="Current feature" value={data?.currentFeature || 'none'} mono />
        <InfoRow label="Server port" value="7300" />
        <InfoRow label="Poll interval" value="30s" last />
      </div>
    </div>
  )
}

function InfoRow({ label, value, mono = false, last = false }: {
  label: string; value: string; mono?: boolean; last?: boolean
}) {
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      gap: 16,
      padding: '12px 20px',
      borderBottom: last ? 'none' : '1px solid rgba(255,255,255,0.06)',
    }}>
      <span style={{ fontSize: 13, color: '#9c9c9d', letterSpacing: '0.2px' }}>{label}</span>
      <span style={{
        fontSize: 13,
        fontWeight: 500,
        color: '#f9f9f9',
        letterSpacing: '0.2px',
        fontFamily: mono ? 'monospace' : 'inherit',
        maxWidth: 300,
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
      }}>
        {value}
      </span>
    </div>
  )
}
